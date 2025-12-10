FROM alpine:3.23

RUN apk add --no-cache \
    restic \
    rsync \
    openssh-server

RUN echo "\
    PermitRootLogin no \
    PasswordAuthentication no \
    PermitEmptyPasswords no \
    AuthenticationMethods publickey \
    " > /etc/ssh/sshd_config.d/50-game-sync.conf

COPY entrypoint.sh /entrypoint.sh
RUN chmod 700 /entrypoint.sh

EXPOSE 22

CMD ["/entrypoint.sh"]
