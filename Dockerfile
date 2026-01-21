FROM alpine:3.23

ARG GAMESYNC_STATE_PATH="bin/gamesync-state"

RUN apk add --no-cache \
    restic \
    rsync \
    openssh-server \
    tzdata \
    bash

COPY server/sshd.conf /etc/ssh/sshd_config.d/50-game-sync.conf

COPY server/entrypoint.sh /entrypoint.sh
RUN chmod 500 /entrypoint.sh

COPY server/restricted-shell /usr/local/bin/restricted-shell
RUN chmod 555 /usr/local/bin/restricted-shell

COPY ${GAMESYNC_STATE_PATH} /usr/local/bin/gamesync-state
RUN chmod 555 /usr/local/bin/gamesync-state

EXPOSE 22

CMD ["/entrypoint.sh"]
