FROM quay.io/giantswarm/helm-chart-testing:v2.4.0 AS ct

RUN pip freeze > /helm-chart-testing-py-requirements.txt

FROM quay.io/giantswarm/golang:1.14.1-alpine3.11 AS golang

FROM quay.io/giantswarm/conftest:v0.18.1 AS conftest

# Build Image
FROM quay.io/giantswarm/alpine:3.11

# Copy go from golang image.
COPY --from=golang /usr/local/go /usr/local/go

# Copy files needed for Helm Chart testing
COPY --from=ct /helm-chart-testing-py-requirements.txt /helm-chart-testing-py-requirements.txt
COPY --from=ct /usr/local/bin/ct /usr/local/bin/ct
COPY --from=ct /etc/ct/chart_schema.yaml /etc/ct/chart_schema.yaml
COPY --from=ct /etc/ct/lintconf.yaml /etc/ct/lintconf.yaml

COPY --from=conftest /usr/local/bin/conftest /usr/local/bin/conftest

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ARG HELM_VERSION=v2.16.3
ARG GOLANGCI_LINT_VERSION=v1.23.8

RUN apk add --no-cache \
        bash \
        ca-certificates \
        curl \
        docker \
        git \
        py-pip \
        openssh-client &&\
        curl -SL https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-linux-amd64.tar.gz | \
            tar -C /usr/bin --strip-components 1 -xvzf - linux-amd64/helm && \
        curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
            sh -s -- -b $GOPATH/bin ${GOLANGCI_LINT_VERSION}

# Setup ssh config for github.com
RUN mkdir ~/.ssh &&\
    chmod 700 ~/.ssh &&\
    ssh-keyscan github.com >> ~/.ssh/known_hosts &&\
    printf "Host github.com\n IdentitiesOnly yes\n IdentityFile ~/.ssh/id_rsa\n" >> ~/.ssh/config &&\
    chmod 600 ~/.ssh/*

RUN pip install -r /helm-chart-testing-py-requirements.txt

ADD ./architect /usr/bin/architect
ENTRYPOINT ["/usr/bin/architect"]
