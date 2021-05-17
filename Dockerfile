FROM golang as build

RUN apt-get update
RUN apt-get install -y libnetfilter-log-dev
RUN go get github.com/goreleaser/goreleaser

ENTRYPOINT [ "goreleaser" ]

#RUN GOOS=linux GOARCH=amd64 go build  -o app github.com/schidstorm/nftables-exporter
