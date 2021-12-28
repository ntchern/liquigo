#!/bin/sh

VERSION=0.1.0

# GOOS=linux GOARCH=amd64 go build -o liquigo-exec
# GOOS=linux GOARCH=arm64 GOARM=7 go build -o liquigo-exec

# docker buildx create --name mybuilder
# docker buildx use mybuilder
docker buildx build \
    --platform linux/amd64,linux/arm64,linux/arm/v7 \
    -t ntchern/liquigo:${VERSION} .

# docker build -t ntchern/liquigo:${VERSION} .
# docker push ntchern/liquigo:${VERSION}
