# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-alpine AS build

WORKDIR /app
ADD ./* ./
RUN go mod download
RUN go build -o /server

##
## Deploy
##

FROM alpine:3.16

WORKDIR /
COPY --from=build /server /server
EXPOSE 8080
USER app:app

ENTRYPOINT ["/server"]