FROM alpine:3.7 AS runtime
RUN apk add --no-cache curl jq bash

FROM golang:1.10.3-alpine3.7 AS build

RUN apk add --no-cache git ;\
  go get -u github.com/golang/dep/cmd/dep

ADD . /go/src/github.com/adgear/git-tags-resource/
RUN  cd /go/src/github.com/adgear/git-tags-resource ;\
  dep ensure ;\
  go build -ldflags "-X main.version=`cat VERSION`" .

# resource-template Dockerfile used to build docker image on concourse ci
FROM runtime
# REQUIRED BY CONCOURSE RESOURCE
ADD check /opt/resource/check
ADD in /opt/resource/in
ADD out /opt/resource/out
COPY --from=build /go/src/github.com/adgear/git-tags-resource/git-tags-resource /usr/local/bin/.

RUN chmod +x /opt/resource/*

WORKDIR /opt/resource