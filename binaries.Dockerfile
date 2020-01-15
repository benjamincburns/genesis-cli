FROM golang:1.13.4-alpine as build

ENV GO111MODULE on
ENV OUTPUT_DIR=/opt/genesis
WORKDIR /go/src/github.com/whiteblock/genesis-cli

RUN apk add git gcc libc-dev make

COPY go.mod go.sum ./
RUN go mod download

COPY . .


