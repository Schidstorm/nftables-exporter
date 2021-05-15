FROM goreleaser/goreleaser:v0.164.0 as build

RUN apk update
RUN apk add libnetfilter_log-dev

RUN mkdir -p /go/src/github.com/schidstorm/nftables-exporter
WORKDIR /go/src/github.com/schidstorm/nftables-exporter
COPY ./* .
COPY .git .git

#RUN GOOS=linux GOARCH=amd64 go build  -o app github.com/schidstorm/nftables-exporter
