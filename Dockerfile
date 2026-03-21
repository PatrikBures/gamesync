ARG ALPINE_VERSION="3.23"

FROM alpine:${ALPINE_VERSION}

ARG GAMESYNC_ADMIN_PATH="bin/gamesync-admin"
ARG GAMESYNC_AUTH_PATH="bin/gamesync-auth"
ARG GAMESYNC_WRAPPER_PATH="bin/gamesync-wrapper"

RUN apk add --no-cache \
    rsync \
    openssh-server \
    tzdata \
    bash

COPY --chmod=555 ${GAMESYNC_ADMIN_PATH}     /usr/local/bin/gamesync-admin
COPY --chmod=555 ${GAMESYNC_AUTH_PATH}      /usr/local/bin/gamesync-auth
COPY --chmod=555 ${GAMESYNC_WRAPPER_PATH}   /usr/local/bin/gamesync-wrapper

COPY server/sshd.conf /etc/ssh/sshd_config.d/50-game-sync.conf
COPY --chmod=500 server/entrypoint.sh /entrypoint.sh

EXPOSE 22

CMD ["/entrypoint.sh"]
