# syntax=docker/dockerfile:1.2
ARG GO_VERSION=1.17

FROM debian:bullseye-slim AS base
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*
# install goup
RUN curl -sSf https://raw.githubusercontent.com/owenthereal/goup/master/install.sh \
      | sh -s -- '--skip-prompt'
ENV PATH=$PATH:/root/.go/bin:/root/.go/current/bin
# install go version
RUN goup install ${GO_VERSION}


FROM base AS tfbuild
WORKDIR /src
RUN --mount=type=bind,src=.,target=. \
    go build -o /go/bin/tf .


FROM base
COPY --from=tfbuild /go/bin/tf /usr/local/bin/tf

ENV INPUT=""
ENV OUTPUT=""
ENV GOPATH="/go"
ENV GOCACHE="/cache"
ENV GOMODCACHE="/cache"

VOLUME /cache
WORKDIR /usr/src/transformer
ENTRYPOINT [ "tf" ]
