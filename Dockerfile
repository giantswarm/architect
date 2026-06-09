FROM gsoci.azurecr.io/giantswarm/helm-chart-testing:v3.14.0 AS ct

FROM gsoci.azurecr.io/giantswarm/app-build-suite:2.1.2 AS abs

FROM gsoci.azurecr.io/giantswarm/golang:1.26.4-alpine3.23 AS golang

FROM gsoci.azurecr.io/giantswarm/conftest:v0.68.2 AS conftest

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
ARG HELM_VERSION=v3.21.0

# renovate: datasource=github-releases depName=kubernetes-sigs/kubebuilder
ARG KUBEBUILDER_VERSION=4.14.0

# renovate: datasource=github-releases depName=sonatype-nexus-community/nancy
ARG NANCY_VERSION=v1.2.0

# renovate: datasource=github-releases depName=sigstore/cosign
ARG COSIGN_VERSION=v3.0.6

# renovate: datasource=github-releases depName=hadolint/hadolint
ARG HADOLINT_VERSION=v2.14.0

# renovate: datasource=github-releases depName=Link-/gh-token
ARG GH_TOKEN_VERSION=v2.0.10

# renovate: datasource=github-releases depName=giantswarm/gitsemver
ARG GITSEMVER_VERSION=v2.0.1

# renovate: datasource=github-releases depName=anchore/syft
ARG SYFT_VERSION=v1.45.1

# renovate: datasource=github-releases depName=oras-project/oras
ARG ORAS_VERSION=v1.3.2

# The `kubeconform` tool is used only when Helm Chart is build and published
# with the `architect` executor, which for majority of the project is not the
# case anymore, for they are build and published with the ABS.
# renovate: datasource=github-releases depName=yannh/kubeconform
ARG KUBECONFORM_VERSION=v0.8.0

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

# Install cosign (used by architect-orb for keyless image/chart/binary signing).
RUN curl -sSL -o /usr/bin/cosign https://github.com/sigstore/cosign/releases/download/${COSIGN_VERSION}/cosign-linux-${TARGETARCH} && \
    chmod +x /usr/bin/cosign

# Install hadolint (Dockerfile linter). Upstream asset names use Linux-x86_64 /
# Linux-arm64 instead of linux-amd64 / linux-arm64, so translate.
RUN case "${TARGETARCH}" in \
    amd64) hadolint_arch=x86_64 ;; \
    arm64) hadolint_arch=arm64 ;; \
    *) echo "Unsupported TARGETARCH=${TARGETARCH} for hadolint"; exit 1 ;; \
    esac && \
    curl -sSL -o /usr/bin/hadolint "https://github.com/hadolint/hadolint/releases/download/${HADOLINT_VERSION}/hadolint-Linux-${hadolint_arch}" && \
    chmod +x /usr/bin/hadolint

# Install kubeconform from the upstream pre-built release tarball.
RUN curl -sSL "https://github.com/yannh/kubeconform/releases/download/${KUBECONFORM_VERSION}/kubeconform-linux-${TARGETARCH}.tar.gz" | \
    tar -C /usr/bin -xzf - kubeconform

# Install gh-token that can generate temporary tokens to authenticate towards Github and use it to access the API
RUN wget --no-verbose https://github.com/Link-/gh-token/releases/download/${GH_TOKEN_VERSION}/linux-${TARGETARCH} -O /usr/bin/gh-token && chmod 700 /usr/bin/gh-token

# Install gitsemver CLI for use in CI scripts running inside the container.
RUN curl -sSL "https://github.com/giantswarm/gitsemver/releases/download/${GITSEMVER_VERSION}/gitsemver-${GITSEMVER_VERSION}-linux-${TARGETARCH}.tar.gz" | \
    tar -C /usr/bin --strip-components 1 -xzf - gitsemver-${GITSEMVER_VERSION}-linux-${TARGETARCH}/gitsemver

# Install syft (SBOM generator). Upstream release tarballs use the version
# without the leading `v` in the asset filename, while the download path uses
# the `v`-prefixed tag.
RUN curl -sSL "https://github.com/anchore/syft/releases/download/${SYFT_VERSION}/syft_${SYFT_VERSION#v}_linux_${TARGETARCH}.tar.gz" | \
    tar -C /usr/bin -xzf - syft

# Install oras (OCI registry client). Like syft, the asset filename uses the
# version without the leading `v` while the download path uses the tag.
RUN curl -sSL "https://github.com/oras-project/oras/releases/download/${ORAS_VERSION}/oras_${ORAS_VERSION#v}_linux_${TARGETARCH}.tar.gz" | \
    tar -C /usr/bin -xzf - oras

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
