# syntax=docker/dockerfile:1

FROM golang:1.21 AS builder

WORKDIR /build

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./
COPY ./internal ./

# Build
RUN go build ./cmd/profiles -o /profiles

FROM alpine

WORKDIR /build

COPY --from=builder /build/profiles /build/profiles

EXPOSE 8001

# Run
CMD ["/build/profiles"]