FROM gsoci.azurecr.io/giantswarm/helm-chart-testing:v3.14.0 AS ct

FROM gsoci.azurecr.io/giantswarm/app-build-suite:1.3.0 AS abs

FROM gsoci.azurecr.io/giantswarm/golang:1.25.5-alpine3.23 AS golang

FROM gsoci.azurecr.io/giantswarm/conftest:v0.65.0 AS conftest

# Build Image
FROM gsoci.azurecr.io/giantswarm/alpine:3.23.0

# Copy go from golang image.
COPY --from=golang /usr/local/go /usr/local/go

# Copy files needed for Helm Chart testing
COPY --from=ct /usr/local/bin/ct /usr/local/bin/ct
COPY --from=abs /abs/resources/ct_schemas/gs_metadata_chart_schema.yaml /etc/ct/chart_schema.yaml
COPY --from=ct /etc/ct/lintconf.yaml /etc/ct/lintconf.yaml

COPY --from=conftest /usr/local/bin/conftest /usr/local/bin/conftest

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

ARG TARGETARCH
ARG TARGETOS

# renovate: datasource=github-releases depName=helm/helm
ARG HELM_VERSION=v3.19.2
# renovate: datasource=github-releases depName=kubernetes-sigs/kubebuilder
ARG KUBEBUILDER_VERSION=3.1.0
# renovate: datasource=github-releases depName=sonatype-nexus-community/nancy
ARG NANCY_VERSION=v1.0.52
# The `kubeconform` tool is used only when Helm Chart is build and published
# with the `architect` executor, which for majority of the project is not the
# case anymore, for they are build and published with the ABS.
# renovate: datasource=github-releases depName=yannh/kubeconform
ARG KUBECONFORM_VERSION=v0.7.0
# The `yamale` tool does not seem to be used anymore, it is still here just in case
# some CI magic somewhere still relies on it.
# renovate: datasource=pypi depName=yamale
ARG CT_YAMALE_VER=6.1.0
# renovate: datasource=pypi depName=yamllint
ARG CT_YAMLLINT_VER=1.37.1

RUN apk add --no-cache \
  bash \
  ca-certificates \
  curl \
  docker \
  git \
  github-cli \
  jq \
  py-pip \
  openssh-client \
  make \
  wget \
  yq &&\
  curl -SL https://get.helm.sh/helm-${HELM_VERSION}-${TARGETOS}-${TARGETARCH}.tar.gz | \
  tar -C /usr/bin --strip-components 1 -xvzf - ${TARGETOS}-${TARGETARCH}/helm && \
  curl -sSfL -o /usr/local/kubebuilder https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${KUBEBUILDER_VERSION}/kubebuilder_${TARGETOS}_${TARGETARCH} && \
  curl -sSL -o /usr/bin/nancy https://github.com/sonatype-nexus-community/nancy/releases/download/${NANCY_VERSION}/nancy-${NANCY_VERSION}-${TARGETOS}-${TARGETARCH} && \
  chmod +x /usr/bin/nancy && \
  go install github.com/yannh/kubeconform/cmd/kubeconform@${KUBECONFORM_VERSION}

# Install gh-token that can generate temporary tokens to authenticate towards Github and use it to access the API
RUN wget https://github.com/Link-/gh-token/releases/download/v2.0.6/${TARGETOS}-${TARGETARCH} -O /usr/bin/gh-token && chmod 700 /usr/bin/gh-token

# Setup ssh config for github.com
RUN mkdir ~/.ssh &&\
  chmod 700 ~/.ssh &&\
  ssh-keyscan github.com >> ~/.ssh/known_hosts &&\
  printf "Host github.com\n IdentitiesOnly yes\n IdentityFile ~/.ssh/id_rsa\n" >> ~/.ssh/config &&\
  chmod 600 ~/.ssh/*

# Allow installing python modules in the global context.
# See https://peps.python.org/pep-0668/
RUN rm -f /usr/lib/python3.11/EXTERNALLY-MANAGED

RUN pip install --break-system-packages yamllint==${CT_YAMLLINT_VER} yamale==${CT_YAMALE_VER}

ADD ./architect-${TARGETOS}-${TARGETARCH} /usr/bin/architect
ENTRYPOINT ["/usr/bin/architect"]
