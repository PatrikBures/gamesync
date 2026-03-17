#!/bin/bash

mkdir -p /etc/ssh/keys
# creates host key if it does not exist
if [ ! -f /etc/ssh/keys/ssh_host_ed25519_key ]; then
    echo "created host key"
    ssh-keygen -A
    mv /etc/ssh/ssh_host_*key* /etc/ssh/keys/
fi

gamesync-admin init dirs
gamesync-admin init migrate

# exec replaces the shell with the process sshd so that pid 1 is sshd and docker can stop it cleanly
# -D do not detach
# -e send error to standard error instead of /var/log, docker does not read /var log
exec /usr/sbin/sshd -D -e
