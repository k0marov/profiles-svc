# syntax=docker/dockerfile:1

FROM golang:1.21 AS builder

WORKDIR /build

COPY . .

RUN go mod download


# Build
RUN GOOS=linux GOARCH=amd64 go build -o /profiles ./cmd/profiles

FROM alpine

WORKDIR /build

COPY --from=builder /build/profiles /build/profiles

EXPOSE 8001

# Run
ENTRYPOINT /build/profiles