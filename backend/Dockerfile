# syntax=docker/dockerfile:1
ARG SKAFFOLD_GO_GCFLAGS

### Build executable
FROM golang:1.21 as builder
WORKDIR /go/src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg go mod download
COPY gql ./gql
COPY model ./model
COPY resolver ./resolver
COPY services ./services
COPY util ./util
COPY web ./web
COPY *.go ./
ENV CGO_ENABLED="0"
ENV GO111MODULE="on"
RUN --mount=type=cache,target=/go/pkg go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o /go/bin/server . && chmod +x /go/bin/server

### Target layer
FROM gcr.io/distroless/base-debian11
COPY --from=builder /go/bin/server /usr/local/bin/server
ENV GOTRACEBACK=single
ENTRYPOINT ["/usr/local/bin/server"]
