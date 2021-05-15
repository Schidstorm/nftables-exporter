FROM goreleaser/goreleaser:v0.163.1 as build

RUN apk update
RUN apk add libnetfilter_log-dev

RUN mkdir -p /usr/local/go/src/github.com/schidstorm/nftables-exporter
WORKDIR /usr/local/go/src/github.com/schidstorm/nftables-exporter
COPY ./* .
COPY .git .git

#RUN GOOS=linux GOARCH=amd64 go build  -o app github.com/schidstorm/nftables-exporter
