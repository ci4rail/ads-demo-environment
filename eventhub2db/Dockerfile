FROM golang:1.16 AS build
WORKDIR /go/src/github.com/ci4rail/
COPY . /go/src/github.com/ci4rail/ads-demo-environment
ENV CGO_ENABLED=0
ENV GOPATH=/go
ARG VERSION=dev

WORKDIR /go/src/github.com/ci4rail/ads-demo-environment/eventhub2db
RUN make test && VERSION=${VERSION} make BIN_DIR=/install

FROM alpine:3.12
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=build /install/eventhub2db /eventhub2db
ENTRYPOINT ["/eventhub2db"]
