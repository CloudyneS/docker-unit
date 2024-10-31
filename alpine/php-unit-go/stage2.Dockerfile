FROM ghcr.io/cloudynes/php-unit-go:stage1-8.1-alpine

COPY goapp /src/goapp

RUN apk add --no-cache gcc musl-dev \
    && cd /src/goapp \
    && go get \
    && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags "-L /usr/src/unit/build/lib -L /usr/lib/unit/modules/" -o /app/goapp main.go \
    && chmod +x /app/goapp

