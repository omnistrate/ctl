# syntax=docker/dockerfile:1

## Build
FROM golang:1.22-alpine AS builder

RUN apk update && apk add git make

WORKDIR /go/src/ctl
COPY ./ .

ARG TARGETOS
ARG TARGETARCH
ARG GIT_USER
ARG GIT_TOKEN
ARG GOPROXY
ARG GOSUMDB
ENV CGO_ENABLED=0
ENV GOPRIVATE=github.com/omnistrate
ENV GOPROXY=${GOPROXY}
ENV GOSUMDB=${GOSUMDB}
ENV GOMAXPROCS=8

RUN echo Configuring access to private repos
RUN echo "machine github.com login ${GIT_USER} password ${GIT_TOKEN}" > ~/.netrc

RUN echo Fetching project dependencies
RUN --mount=type=cache,target=/go/pkg/mod \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go mod tidy

RUN --mount=type=cache,target=/go/pkg/mod \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go mod vendor

RUN echo Building and installing Monitoring Service
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make build
RUN echo Short testing to validate the build
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go test ./... -short
## Deploy
FROM alpine:latest AS app

ARG TARGETOS
ARG TARGETARCH

RUN apk update && apk add --no-cache ca-certificates curl

COPY --from=builder /go/src/ctl/dist/omnistrate-ctl-${TARGETOS}-${TARGETARCH} /usr/local/bin/omnistrate-ctl
WORKDIR /omnistrate
ENTRYPOINT [ "/usr/local/bin/omnistrate-ctl" ]