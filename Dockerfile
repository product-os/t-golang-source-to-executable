# syntax=docker/dockerfile:1.2
ARG GO_VERSION=1.17
ARG GOTESTSUM_COMMIT=v0.3.5

FROM docker.io/library/golang:${GO_VERSION} AS base
RUN apt-get update && apt-get install -y \
      curl \
    && rm -rf /var/lib/apt/lists/*


FROM base AS gotestsum
# install gotestsum
ARG GOTESTSUM_COMMIT
RUN CGO_ENABLED=0 go install -buildmode=pie "gotest.tools/gotestsum@${GOTESTSUM_COMMIT}"


FROM base AS gobuild
# compile transformer
WORKDIR /src
RUN --mount=type=bind,src=.,target=. \
    go build -o /go/bin/tf .


FROM base AS final
COPY --from=gotestsum /go/bin/gotestsum /usr/local/bin/
COPY --from=gobuild /go/bin/tf /usr/local/bin/

ENV INPUT=""
ENV OUTPUT=""

ENV GOCACHE="/cache"
ENV GOMODCACHE="/cache"

VOLUME /cache
WORKDIR /usr/src/transformer

ENTRYPOINT [ "tf", "-mode", "build" ]
