FROM golang:1-alpine

RUN apk add protobuf

RUN apk add bash

WORKDIR /usr/app

COPY . .

ENTRYPOINT [ "go", "run", "./cmd/server/server.go" ]




