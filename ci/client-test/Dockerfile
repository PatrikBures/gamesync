ARG ALPINE_VERSION="3.23"

FROM alpine:${ALPINE_VERSION}

RUN apk add --no-cache \
    rsync \
    openssh-client

COPY ./gamesync /usr/local/bin/gamesync
