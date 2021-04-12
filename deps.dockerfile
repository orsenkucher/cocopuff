FROM golang:1.16.3-buster AS deps
WORKDIR /deps
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY pub pub
