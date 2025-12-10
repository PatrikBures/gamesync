FROM alpine:3.23

RUN apk add --no-cache \
    restic \
    rsync \
    openssh-server

COPY entrypoint.sh /entrypoint.sh
RUN chmod 700 /entrypoint.sh

EXPOSE 22

CMD ["/entrypoint.sh"]
