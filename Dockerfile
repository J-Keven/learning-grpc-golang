FROM golang:1-alpine

RUN apk add protobuf

RUN apk add bash

WORKDIR /usr/app


COPY . .

CMD [ "tail", "-f", "/dev/null" ]




