FROM golang:1.21.5-alpine AS base

WORKDIR /var/code
COPY ./ ./

RUN \
    CGO_ENABLED=0 GOOS=linux go build \
    -o /usr/local/bin/server \
    ./cmd/run-server/main.go

FROM alpine:3.16.0 AS app

COPY --from=base /usr/local/bin/server /usr/local/bin/server
ENTRYPOINT ["/usr/local/bin/server"]
