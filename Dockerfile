# syntax = docker/dockerfile:1
ARG GO_BASE_VERSION=1.22
FROM golang:$GO_BASE_VERSION as builder
ARG TARGETOS
ARG TARGETARCH

# TODO remove this line once the repo is public
ENV GOPROXY "https://goproxy.build.ue1.snooguts.net"
ENV GOPRIVATE ""
ENV GONOSUMDB "github.snooguts.net/*"
RUN git config --global url.git@github.snooguts.net:.insteadof https://github.snooguts.net/

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/main.go cmd/main.go
COPY api/ api/
COPY internal/ internal/

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot AS runner-base
WORKDIR /
USER 65532:65532

ENTRYPOINT ["/manager"]

FROM builder as prod-builder
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -v -a -o /manager cmd/main.go

FROM runner-base as production
COPY --from=prod-builder /manager /manager
