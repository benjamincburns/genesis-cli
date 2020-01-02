FROM golang:1.13.4-alpine as build

ENV GO111MODULE on
ENV OUTPUT_DIR=/usr/bin
WORKDIR /go/src/github.com/whiteblock/genesis-cli

RUN apk add git gcc libc-dev make

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make -e

FROM alpine:3.10 as final

RUN apk add ca-certificates

COPY --from=build /usr/bin/genesis /usr/bin/genesis

ENTRYPOINT ["/usr/bin/genesis"]