ARG GOLANG_IMAGE=""

FROM ${GOLANG_IMAGE}

RUN apk add --no-cache \
    rsync \
    openssh-client

COPY ./gamesync /usr/local/bin/gamesync
