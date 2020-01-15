FROM golang:1.13.4-buster as build

ENV GO111MODULE on
ENV OUTPUT_DIR=/opt/genesis
WORKDIR /go/src/github.com/whiteblock/genesis-cli

RUN apt-get install -y git gcc libc-dev make

COPY go.mod go.sum ./
RUN go mod download

COPY . .


