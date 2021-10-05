# syntax=docker/dockerfile:1.2
FROM golang:1.17-alpine

ENV INPUT="" \
    OUTPUT=""

RUN apk add --no-cache bash jq yq

WORKDIR /usr/src/transformer

COPY . ./
ENTRYPOINT [ "./scripts/entrypoint.sh" ]
