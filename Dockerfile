# Stage 0
FROM quay.io/giantswarm/golang:1.13.0-alpine3.10

# Stage 1
FROM quay.io/giantswarm/alpine:3.10

# Copy go from golang image.
COPY --from=0 /usr/local/go /usr/local/go

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ARG HELM_VERSION=v2.14.0

RUN apk add --no-cache \
        bash \
        ca-certificates \
        curl \
        docker \
        git \
        openssh-client &&\
        curl -SL https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-linux-amd64.tar.gz | \
            tar -C /usr/bin --strip-components 1 -xvzf - linux-amd64/helm

# Setup ssh config for github.com
RUN mkdir ~/.ssh &&\
    chmod 700 ~/.ssh &&\
    ssh-keyscan github.com >> ~/.ssh/known_hosts &&\
    printf "Host github.com\n IdentitiesOnly yes\n IdentityFile ~/.ssh/id_rsa\n" >> ~/.ssh/config &&\
    chmod 600 ~/.ssh/*

ADD ./architect /usr/bin/architect
ENTRYPOINT ["/usr/bin/architect"]
