FROM alpine:3.23

RUN apk add --no-cache \
    restic \
    rsync \
    openssh-server

COPY sshd.conf /etc/ssh/sshd_config.d/50-game-sync.conf

COPY entrypoint.sh /entrypoint.sh
RUN chmod 500 /entrypoint.sh

EXPOSE 22

CMD ["/entrypoint.sh"]
