#!/bin/sh

VERSION=0.1.0

GOOS=linux GOARCH=amd64 go build -o liquigo-exec

docker build -t ntchern/liquigo:${VERSION} .
docker push ntchern/liquigo:${VERSION}
