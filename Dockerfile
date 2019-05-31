FROM alpine:3.9

ARG HELM_VERSION=v2.14.0

RUN apk add --no-cache \
        ca-certificates \
        curl \
        git \
        openssh-client &&\
        # Install helm
        curl -SL https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-linux-amd64.tar.gz | \
                tar -C /usr/bin --strip-components 1 -xvzf - linux-amd64/helm &&\
        apk del curl &&\
        rm -f /var/cache/apk/*

# Setup ssh config for github.com
RUN mkdir ~/.ssh &&\
    chmod 700 ~/.ssh &&\
    ssh-keyscan github.com >> ~/.ssh/known_hosts &&\
    printf "Host github.com\n IdentitiesOnly yes\n IdentityFile ~/.ssh/id_rsa\n" >> ~/.ssh/config &&\
    chmod 600 ~/.ssh/*

ADD ./architect /usr/bin/architect
ENTRYPOINT ["/usr/bin/architect"]
