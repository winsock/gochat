# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-bullseye AS build

WORKDIR /app
ADD . ./
RUN go mod download
RUN go build -ldflags="-extldflags=-static" -o /server

##
## Deploy
##

FROM debian:bullseye-slim

WORKDIR /
COPY --from=build /server /server
EXPOSE 8080

ENTRYPOINT ["/server"]