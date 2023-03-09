# syntax=docker/dockerfile:1
ARG SKAFFOLD_GO_GCFLAGS

### Build executable
FROM golang:1.19 as builder
WORKDIR /workspace
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg go mod download
COPY admin ./admin
COPY auth ./auth
COPY operations ./operations
COPY public ./public
COPY util ./util
COPY *.go ./
ENV CGO_ENABLED="0"
ENV GO111MODULE="on"
RUN --mount=type=cache,target=/go/pkg go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o server . && chmod +x server

### Target layer
FROM gcr.io/distroless/base-debian11
COPY --from=builder /workspace/server /usr/local/bin/api
ENV GOTRACEBACK=single
ENTRYPOINT ["/usr/local/bin/api"]
