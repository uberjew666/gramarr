FROM golang:1.17-alpine as builder

RUN mkdir -p /go/src/github.com/alcomoraes/gramarr
WORKDIR /go/src/github.com/alcmoraes/gramarr

RUN apk --update upgrade \
    && apk --no-cache --no-progress add git \
    && rm -rf /var/cache/apk/

ADD . /go/src/github.com/alcmoraes/gramarr/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -installsuffix nocgo -o /gramarr

FROM alpine:3.14
RUN apk --update upgrade \
    && apk --no-cache --no-progress add procps \
    && rm -rf /var/cache/apk/*

COPY --from=builder /gramarr ./

VOLUME ["/config"]

COPY config.json.template /config/config.json

ENTRYPOINT ["/gramarr", "-configDir=/config"]
