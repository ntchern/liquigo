#!/bin/sh

VERSION=0.1.0

docker run --rm --network host \
    -v $(pwd):/workdir \
    -w /workdir \
    -e DB_HOST=localhost \
    -e DB_PORT=5454 \
    -e DB_DATABASE=postgres \
    -e DB_SCHEMA=public \
    -e DB_USER=postgres \
    -e DB_PASSWORD=postgres \
    -e DB_SSLMODE=disable \
    ntchern/liquigo:0.1.0 \
    update \
    --url qq \
    --changeLog /workdir/test-files/2-initial-update/_changelog.yaml
