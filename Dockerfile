FROM golang:1.9

ENV GOPATH=/app
ENV PATH=$PATH:$GOPATH/bin
ENV SERVER_HOSTNAME=gochat_server

WORKDIR /app
ADD $PWD /app/src/github.com/shavit/gochat
