FROM golang:alpine AS builder

ADD . /src
WORKDIR /src
RUN go build cmd/apcupsd_exporter/main.go

FROM alpine:latest

LABEL org.opencontainers.image.source https://github.com/fire833/apcupsd_exporter

RUN apk update
RUN apk upgrade

COPY --from=builder /src/main /apcupsd_exporter

EXPOSE 9162

ENTRYPOINT ["/apcupsd_exporter"]