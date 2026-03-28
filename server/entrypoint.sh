#!/bin/sh

mkdir -p /etc/ssh/keys
# creates host key if it does not exist
if [ ! -f /etc/ssh/keys/ssh_host_ed25519_key ]; then
    echo "created host key"
    ssh-keygen -A || return $?
    mv /etc/ssh/ssh_host_*key* /etc/ssh/keys/
fi

if ! id -u gamesync > /dev/null 2>&1; then
    adduser -H -D -u 1000 -h /data/save -s /bin/sh gamesync || return $?
    passwd -u gamesync || return $?
fi

gamesync-admin init dirs    || return 20
gamesync-admin init migrate || return 21
gamesync-admin init perms   || return 22
gamesync-admin init roles   || return 23

# exec replaces the shell with the process sshd so that pid 1 is sshd and docker can stop it cleanly
# -D do not detach
# -e send error to standard error instead of /var/log, docker does not read /var log
exec /usr/sbin/sshd -D -e
