# TODO:

# Make 3-stages build:
#   1. Load dependencies from go.mod
#   2. Build server binary
#   3. Transfer binaries to buster-slim

FROM golang:1.16.3-buster AS deps
WORKDIR /deps
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM deps AS build
COPY . .
