FROM golang:latest as builder
ENV GOPATH=/go
ENV GO111MODULE=on
WORKDIR ${GOPATH}/src/github.com/Shikugawa/pcp
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o ./main

FROM ubuntu:latest
WORKDIR /app

# service binary
COPY --from=builder /go/src/github.com/Shikugawa/pcp/main .
RUN chmod +x ./main
