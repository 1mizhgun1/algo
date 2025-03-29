FROM golang:1.23.2-alpine AS builder

COPY . /algo/
WORKDIR /algo/

RUN go clean --modcache
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./main.go

ENV TZ="Europe/Moscow"
ENV ZONEINFO=/zoneinfo.zip

EXPOSE 8080

ENTRYPOINT ["./.bin"]
