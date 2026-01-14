FROM alpine:3.23

RUN apk add --no-cache \
    restic \
    rsync \
    openssh-server \
    tzdata

COPY sshd.conf /etc/ssh/sshd_config.d/50-game-sync.conf

COPY entrypoint.sh /entrypoint.sh
RUN chmod 500 /entrypoint.sh

COPY restricted-shell /usr/local/bin/restricted-shell
RUN chmod 555 /usr/local/bin/restricted-shell

COPY gamesync-state /usr/local/bin/gamesync-state
RUN chmod 555 /usr/local/bin/gamesync-state

EXPOSE 22

CMD ["/entrypoint.sh"]
