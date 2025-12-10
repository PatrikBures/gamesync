#!/bin/sh

# creates host key if it does not exist
if [ ! -f /etc/ssh/ssh_host_ed25519_key ]; then
    echo "created host key"
    ssh-keygen -A
fi

chmod 700 /root/.ssh
chmod 600 /root/.ssh/authorized_keys

# exec replaces the shell with the process sshd so that pid 1 is sshd and docker can stop it cleanly
# -D do not detach
# -e send error to standard error instead of /var/log, docker does not read /var log
exec /usr/sbin/sshd -D -e
