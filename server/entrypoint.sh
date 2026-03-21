#!/bin/bash

mkdir -p /etc/ssh/keys
# creates host key if it does not exist
if [ ! -f /etc/ssh/keys/ssh_host_ed25519_key ]; then
    echo "created host key"
    ssh-keygen -A
    mv /etc/ssh/ssh_host_*key* /etc/ssh/keys/
fi

if ! id -u gamesync &> /dev/null; then
    adduser gamesync -D -s /bin/sh
    passwd -u gamesync
fi
addgroup db
adduser gamesync db

gamesync-admin init dirs
gamesync-admin init migrate
gamesync-admin init perms
gamesync-admin init roles

chown -R :db /data/db
chmod 770 /data/db
chmod 660 /data/db/*

# exec replaces the shell with the process sshd so that pid 1 is sshd and docker can stop it cleanly
# -D do not detach
# -e send error to standard error instead of /var/log, docker does not read /var log
exec /usr/sbin/sshd -D -e
