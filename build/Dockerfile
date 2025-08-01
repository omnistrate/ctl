# syntax=docker/dockerfile:1

## Build
FROM golang:1.24-alpine AS builder

RUN apk update && apk add git make

WORKDIR /go/src/ctl
COPY ./ .

ARG TARGETOS
ARG TARGETARCH
ARG GOPROXY
ARG GOSUMDB
ARG GIT_COMMIT
ARG GIT_VERSION
ENV CGO_ENABLED=0
ENV GIT_COMMIT=${GIT_COMMIT}
ENV GIT_VERSION=${GIT_VERSION}

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

RUN apk update && apk add --no-cache ca-certificates curl jq

COPY --from=builder /go/src/ctl/dist/omnistrate-ctl-${TARGETOS}-${TARGETARCH} /usr/local/bin/omnistrate-ctl
WORKDIR /omnistrate
ENTRYPOINT [ "/usr/local/bin/omnistrate-ctl" ]
