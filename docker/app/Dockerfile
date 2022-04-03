# syntax=docker/dockerfile:1

FROM golang:1.16-alpine as base
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build the binary
FROM base AS builder
RUN CGO_ENABLED=0 GOARCH=amd64 go build -o /main cmd/app/main.go   
# Build minimal image
FROM alpine:3
COPY --from=builder main /bin/main
ENTRYPOINT ["/bin/main"]