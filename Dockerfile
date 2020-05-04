FROM golang:1.13 AS build

RUN mkdir -p /go/src/github.com/alcmoraes/gramarr

WORKDIR /go/src/github.com/alcmoraes/gramarr

COPY . .

RUN go get

RUN mkdir -p /app \
             /config

RUN go build -o /app/gramarr

COPY config.json /config/config.json

CMD ["/app/gramarr", "-configDir=/config"]

VOLUME /config