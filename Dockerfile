FROM goreleaser/goreleaser:v0.163.1 as build

ARG GO_MODULE_NAME=github.com/schidstorm/nftables-exporter

RUN apk update
RUN apk add libnetfilter_log-dev

ENV GO111MODULE=on
RUN mkdir -p /go/src/$GO_MODULE_NAME
COPY . /go/src/$GO_MODULE_NAME
WORKDIR /go/src/$GO_MODULE_NAME

#RUN GOOS=linux GOARCH=amd64 go build  -o app github.com/schidstorm/nftables-exporter
