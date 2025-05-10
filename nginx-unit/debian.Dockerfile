ARG SOURCE_VERSION=8.2
FROM golang:1.24.3-bookworm AS goget
FROM ghcr.io/cloudynes/php-extended:${SOURCE_VERSION}-debian

COPY --from=goget /usr/local/go /usr/local/go
COPY --from=goget /go /go

ENV GOTOOLCHAIN=local
ENV GOBIN=/usr/local/go/bin
ENV GOPATH=/go
ENV GOLANG_VERSION=1.24.3