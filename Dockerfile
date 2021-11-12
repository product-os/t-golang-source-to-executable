# syntax=docker/dockerfile:1.2
ARG GO_VERSION=1.17
ARG GOTESTSUM_COMMIT=v0.3.5

FROM debian:bullseye-slim AS base
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# install goup
RUN curl -sSf https://raw.githubusercontent.com/owenthereal/goup/master/install.sh \
      | sh -s -- '--skip-prompt'
ENV PATH=$PATH:/root/.go/bin:/root/.go/current/bin:/root/go/bin

# install go version
ARG GO_VERSION
RUN goup install ${GO_VERSION}

# install gotestsum
ARG GOTESTSUM_COMMIT
RUN CGO_ENABLED=0 go install -buildmode=pie "gotest.tools/gotestsum@${GOTESTSUM_COMMIT}"

# compile transformer
WORKDIR /src
RUN --mount=type=bind,src=.,target=. \
    go build -o /usr/local/bin/tf .

ENV INPUT=""
ENV OUTPUT=""

ENV GOCACHE="/cache"
ENV GOMODCACHE="/cache"

VOLUME /cache
WORKDIR /usr/src/transformer

FROM base AS test
ENTRYPOINT [ "tf", "-mode", "test" ]

FROM base AS build
ENTRYPOINT [ "tf", "-mode", "build" ]
