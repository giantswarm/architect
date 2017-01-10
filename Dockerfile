FROM ubuntu:16.04
# LTS - Xenial Xerus

# install packages we need to add ppas
RUN apt-get -y update
RUN apt-get install -y \
    apt-transport-https=1.2.15ubuntu0.2 \ 
    software-properties-common=0.96.20.5

# add docker repository
RUN apt-key adv \
    --keyserver hkp://p80.pool.sks-keyservers.net:80 \
    --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
RUN echo "deb https://apt.dockerproject.org/repo ubuntu-xenial main" | tee /etc/apt/sources.list.d/docker.list

# add the glide ppa
RUN add-apt-repository ppa:masterminds/glide
RUN apt-get -y update

# install packages we need (wget to install golang, docker, git, and glide for builds)
RUN apt-get install -y \
    docker-engine=1.11.1-0~xenial \
    git=1:2.7.4-0ubuntu1 \
    glide=0.12.3~xenial \
    wget=1.17.1-1ubuntu1.1

# install golang
RUN wget https://storage.googleapis.com/golang/go1.7.3.linux-amd64.tar.gz -qO- | tar xzf - \
    && mv ./go /usr/local/
    
ADD ./architect /architect

ENTRYPOINT ["/architect"]
