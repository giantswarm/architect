FROM quay.io/giantswarm/helm-chart-testing:v3.6.0 AS ct

FROM quay.io/giantswarm/app-build-suite:0.2.3 AS abs

FROM quay.io/giantswarm/golang:1.18.1-alpine3.15 AS golang

FROM quay.io/giantswarm/conftest:v0.33.1 AS conftest

# Build Image
FROM quay.io/giantswarm/alpine:3.14.2

# Copy go from golang image.
COPY --from=golang /usr/local/go /usr/local/go

# Copy files needed for Helm Chart testing
COPY --from=ct /usr/local/bin/ct /usr/local/bin/ct
COPY --from=abs /abs/resources/ct_schemas/gs_metadata_chart_schema.yaml /etc/ct/chart_schema.yaml
COPY --from=ct /etc/ct/lintconf.yaml /etc/ct/lintconf.yaml

COPY --from=conftest /usr/local/bin/conftest /usr/local/bin/conftest

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ARG HELM_VERSION=v3.8.1
ARG KUBEBUILDER_VERSION=3.1.0
ARG GOLANGCI_LINT_VERSION=v1.46.0
ARG NANCY_VERSION=v1.0.37
ARG KUBEVAL_VERSION=v0.16.1
ARG CT_YAMALE_VER=3.0.6
ARG CT_YAMLLINT_VER=1.26.1

RUN apk add --no-cache \
        bash \
        ca-certificates \
        curl \
        docker \
        git \
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
        curl -sSL https://github.com/instrumenta/kubeval/releases/download/${KUBEVAL_VERSION}/kubeval-$(go env GOOS)-$(go env GOARCH).tar.gz | \
            tar -C /usr/bin -xvzf - kubeval

# Setup ssh config for github.com
RUN mkdir ~/.ssh &&\
    chmod 700 ~/.ssh &&\
    ssh-keyscan github.com >> ~/.ssh/known_hosts &&\
    printf "Host github.com\n IdentitiesOnly yes\n IdentityFile ~/.ssh/id_rsa\n" >> ~/.ssh/config &&\
    chmod 600 ~/.ssh/*

RUN pip install yamllint==${CT_YAMLLINT_VER} yamale==${CT_YAMALE_VER}

ADD ./architect /usr/bin/architect
ENTRYPOINT ["/usr/bin/architect"]
