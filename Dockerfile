#!/usr/bin/env bash

ARG GO_VERSION=1.13
ARG ALPINE_VERSION=3.10

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build-stage
ENV CGO_ENABLED 0
ENV GO111MODULE=on

ADD . /usr/src/app
WORKDIR /usr/src/app

RUN set -ex && \
    apk update && \
    apk add --no-cache git && \
    apk add linux-headers && \
    apk add musl && \
    apk add zlib-dev && \
    apk add libjpeg-turbo-dev && \
    apk add gcc && \
    apk add g++ && \
    apk add make && \
    apk add gfortran && \
    apk add openblas-dev && \
    apk add zlib-dev

RUN go get github.com/pilu/fresh && \
    go get github.com/RyanCarrier/dijkstra && \
    go get github.com/gin-gonic/gin && \
    go get github.com/gin-contrib/cors && \
    go get github.com/twpayne/go-polyline && \
    go get github.com/go-sql-driver/mysql
