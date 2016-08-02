FROM ubuntu:14.04

MAINTAINER tatva@lindenlab.com

### Update system
RUN apt-get update && apt-get -y upgrade
RUN apt-get install -y wget ca-certificates build-essential git mercurial bzr

RUN echo "Installing 1.5"

# install golang 1.5

RUN wget --no-verbose https://storage.googleapis.com/golang/go1.5.4.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.5.4.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin
ENV GOPATH /usr/local/go/


RUN apt-get install -y cmake python-sphinx protobuf-compiler debhelper && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Get heka
RUN mkdir /heka
RUN git clone https://github.com/mozilla-services/heka.git /heka
COPY plugin_loader.cmake /heka/cmake/plugin_loader.cmake

# get oauth2 & other dependencies
RUN mkdir -p /heka/build/heka/src
WORKDIR /heka/build/heka/src

RUN go get golang.org/x/oauth2
RUN go get golang.org/x/oauth2/google
RUN go get golang.org/x/oauth2/jwt
RUN go get google.golang.org/api/bigquery/v2
RUN go get golang.org/x/net/context
RUN go get google.golang.org/cloud/pubsub
RUN go get github.com/aranair/heka-bigquery/bq

WORKDIR /heka
RUN git checkout versions/0.10
RUN /bin/bash -c 'source build.sh'

COPY heka-config.toml.tmpl /app/heka-config.toml.tmpl

RUN mkdir -p /usr/share/heka/lua_decoders && \
    mkdir -p /usr/share/heka/lua_encoders && \
    mkdir -p /usr/share/heka/lua_filters && \
    mkdir -p /usr/share/heka/lua_modules && \
    cp /heka/sandbox/lua/decoders/* /usr/share/heka/lua_decoders && \
    cp /heka/sandbox/lua/encoders/* /usr/share/heka/lua_encoders && \
    cp /heka/sandbox/lua/filters/* /usr/share/heka/lua_filters && \
    cp /heka/build/heka/lib/luasandbox/modules/* /usr/share/heka/lua_modules

EXPOSE 4881

WORKDIR /heka/build/heka/bin
#CMD ["/heka/build/heka/bin/hekad", "--config", "/app/config.toml"]