# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:alpine AS build

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 go build -o /liquigo

##
## Deploy
##
FROM alpine

COPY --from=build /liquigo /

ENTRYPOINT ["/liquigo"]
