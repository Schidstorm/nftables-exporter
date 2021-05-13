FROM golang as build

COPY . /code
WORKDIR /code
RUN apt update
RUN apt-get install -y libnetfilter-log-dev
RUN GOOS=linux GOARCH=amd64 go build  -o app -a .

FROM scratch
COPY --from=build /code/app /