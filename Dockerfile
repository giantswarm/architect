FROM gsoci.azurecr.io/giantswarm/helm-chart-testing:v3.14.0 AS ct

FROM gsoci.azurecr.io/giantswarm/app-build-suite:1.8.0 AS abs

FROM gsoci.azurecr.io/giantswarm/golang:1.26.3-alpine3.23 AS golang

FROM gsoci.azurecr.io/giantswarm/conftest:v0.68.2 AS conftest

# Cross-compile kubeconform on the build host instead of emulating the target
# arch under QEMU when buildx targets a non-host platform.
FROM --platform=$BUILDPLATFORM gsoci.azurecr.io/giantswarm/golang:1.26.2-alpine3.23 AS kubeconform-builder
ARG TARGETOS
ARG TARGETARCH
# renovate: datasource=github-releases depName=yannh/kubeconform
ARG KUBECONFORM_VERSION=v0.7.0
ENV GOPATH=/go
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go install github.com/yannh/kubeconform/cmd/kubeconform@${KUBECONFORM_VERSION}
# `go install` puts the binary at /go/bin/<name> when GOOS/GOARCH match the
# host, and at /go/bin/${GOOS}_${GOARCH}/<name> when cross-compiling.
# Normalize so the final stage has a single known source path.
RUN mkdir -p /out && \
    if [ -f "/go/bin/kubeconform" ]; then \
      cp /go/bin/kubeconform /out/kubeconform; \
    else \
      cp "/go/bin/${TARGETOS}_${TARGETARCH}/kubeconform" /out/kubeconform; \
    fi

# Build Image
FROM gsoci.azurecr.io/giantswarm/alpine:3.23.4

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

# renovate: datasource=github-releases depName=helm/helm
ARG HELM_VERSION=v3.20.2

# renovate: datasource=github-releases depName=kubernetes-sigs/kubebuilder
ARG KUBEBUILDER_VERSION=4.14.0

# renovate: datasource=github-releases depName=sonatype-nexus-community/nancy
ARG NANCY_VERSION=v1.2.0

# The `yamale` tool does not seem to be used anymore, it is still here just in case
# some CI magic somewhere still relies on it.
# renovate: datasource=pypi depName=yamale
ARG CT_YAMALE_VER=6.1.0

# renovate: datasource=pypi depName=yamllint
ARG CT_YAMLLINT_VER=1.38.0

RUN apk add --no-cache --no-scripts \
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
  yq

SHELL ["/bin/bash", "-xc"]

# Install Helm
RUN curl -sSL https://get.helm.sh/helm-${HELM_VERSION}-linux-${TARGETARCH}.tar.gz | \
  tar -C /usr/bin --strip-components 1 -xvzf - linux-${TARGETARCH}/helm

# Install kubebuilder
RUN curl -sSfL -o /usr/local/kubebuilder https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${KUBEBUILDER_VERSION}/kubebuilder_linux_${TARGETARCH}

# Install nancy
RUN curl -sSL -o /usr/bin/nancy https://github.com/sonatype-nexus-community/nancy/releases/download/${NANCY_VERSION}/nancy-${NANCY_VERSION}-linux-${TARGETARCH} && \
  chmod +x /usr/bin/nancy

# Install kubeconform (pre-built on the host arch in the kubeconform-builder stage above).
# Used only when Helm charts are built and published with the `architect` executor,
# which most projects no longer do (they use ABS instead).
COPY --from=kubeconform-builder /out/kubeconform /go/bin/kubeconform

# Install gh-token that can generate temporary tokens to authenticate towards Github and use it to access the API
RUN wget --no-verbose https://github.com/Link-/gh-token/releases/download/v2.0.6/linux-${TARGETARCH} -O /usr/bin/gh-token && chmod 700 /usr/bin/gh-token

# Setup ssh config for github.com
RUN mkdir ~/.ssh && \
  chmod 700 ~/.ssh && \
  ssh-keyscan github.com >> ~/.ssh/known_hosts &&\
  printf "Host github.com\n IdentitiesOnly yes\n IdentityFile ~/.ssh/id_rsa\n" >> ~/.ssh/config && \
  chmod 600 ~/.ssh/known_hosts ~/.ssh/config

# Allow installing python modules in the global context.
# See https://peps.python.org/pep-0668/
RUN rm -f /usr/lib/python3.11/EXTERNALLY-MANAGED

RUN pip install --break-system-packages --root-user-action=ignore yamllint==${CT_YAMLLINT_VER} yamale==${CT_YAMALE_VER}

ADD ./architect-linux-${TARGETARCH} /usr/bin/architect
ENTRYPOINT ["/usr/bin/architect"]
