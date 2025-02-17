FROM gsoci.azurecr.io/giantswarm/helm-chart-testing:v3.12.0 AS ct

FROM gsoci.azurecr.io/giantswarm/app-build-suite:1.2.8 AS abs

FROM gsoci.azurecr.io/giantswarm/golang:1.23.4-alpine3.19 AS golang

FROM gsoci.azurecr.io/giantswarm/conftest:v0.57.0 AS conftest

# Build Image
FROM gsoci.azurecr.io/giantswarm/alpine:3.19

# Copy go from golang image.
COPY --from=golang /usr/local/go /usr/local/go

# Copy files needed for Helm Chart testing
COPY --from=ct /usr/local/bin/ct /usr/local/bin/ct
COPY --from=abs /abs/resources/ct_schemas/gs_metadata_chart_schema.yaml /etc/ct/chart_schema.yaml
COPY --from=ct /etc/ct/lintconf.yaml /etc/ct/lintconf.yaml

COPY --from=conftest /usr/local/bin/conftest /usr/local/bin/conftest

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# renovate: datasource=github-releases depName=helm/helm
ARG HELM_VERSION=v3.17.1
# renovate: datasource=github-releases depName=kubernetes-sigs/kubebuilder
ARG KUBEBUILDER_VERSION=3.1.0
# renovate: datasource=github-releases depName=golangci/golangci-lint
ARG GOLANGCI_LINT_VERSION=v1.64.5
# renovate: datasource=github-releases depName=sonatype-nexus-community/nancy
ARG NANCY_VERSION=v1.0.46
# renovate: datasource=github-releases depName=yannh/kubeconform
ARG KUBECONFORM_VERSION=v0.4.14
# renovate: datasource=pypi depName=yamale
ARG CT_YAMALE_VER=3.0.8
# renovate: datasource=pypi depName=yamllint
ARG CT_YAMLLINT_VER=1.35.1

RUN apk add --no-cache \
  bash \
  ca-certificates \
  curl \
  docker \
  git \
  jq \
  py-pip \
  openssh-client \
  make \
  yq &&\
  curl -SL https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz | \
  tar -C /usr/bin --strip-components 1 -xvzf - linux-amd64/helm && \
  curl -sSfL -o /usr/local/kubebuilder https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${KUBEBUILDER_VERSION}/kubebuilder_$(go env GOOS)_$(go env GOARCH) && \
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
  sh -s -- -b $GOPATH/bin ${GOLANGCI_LINT_VERSION} && \
  curl -sSL -o /usr/bin/nancy https://github.com/sonatype-nexus-community/nancy/releases/download/${NANCY_VERSION}/nancy-${NANCY_VERSION}-linux-amd64 && \
  chmod +x /usr/bin/nancy && \
  go install github.com/yannh/kubeconform/cmd/kubeconform@${KUBECONFORM_VERSION}

# Setup ssh config for github.com
RUN mkdir ~/.ssh &&\
  chmod 700 ~/.ssh &&\
  ssh-keyscan github.com >> ~/.ssh/known_hosts &&\
  printf "Host github.com\n IdentitiesOnly yes\n IdentityFile ~/.ssh/id_rsa\n" >> ~/.ssh/config &&\
  chmod 600 ~/.ssh/*

# Allow installing python modules in the global context.
# See https://peps.python.org/pep-0668/
RUN rm -f /usr/lib/python3.11/EXTERNALLY-MANAGED

RUN pip install yamllint==${CT_YAMLLINT_VER} yamale==${CT_YAMALE_VER}

ADD ./architect /usr/bin/architect
ENTRYPOINT ["/usr/bin/architect"]
