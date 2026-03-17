ARG ALPINE_VERSION="3.23"

FROM alpine:${ALPINE_VERSION}

ARG GAMESYNC_STATE_PATH="bin/gamesync-state"
ARG GAMESYNC_ADMIN_PATH="bin/gamesync-admin"

RUN apk add --no-cache \
    restic \
    rsync \
    openssh-server \
    tzdata \
    bash

COPY server/sshd.conf /etc/ssh/sshd_config.d/50-game-sync.conf
COPY --chmod=500 server/entrypoint.sh /entrypoint.sh
COPY --chmod=555 server/restricted-shell /usr/local/bin/restricted-shell
COPY --chmod=555 ${GAMESYNC_STATE_PATH} /usr/local/bin/gamesync-state
COPY --chmod=555 ${GAMESYNC_ADMIN_PATH} /usr/local/bin/gamesync-admin

EXPOSE 22

CMD ["/entrypoint.sh"]
