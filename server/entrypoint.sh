#!/bin/sh

REPO_DIR=/data/repos
SAVE_DIR=/data/saves

create_users() {
    for file in /config/users/*; do
        if [ "$file" = "*" ]; then
            continue
        fi

        user="$(basename "${file}")"

        # uses id to check if user exists,
        # the -u option is used because it prints out less info making it faster
        if id -u "$user" > /dev/null 2>&1 ; then
            echo "$user already exists"
            continue
        fi

        # -D no password
        adduser "$user" -D > /dev/null
        passwd -u "$user"


        # adds ssh key to authorized_keys
        ssh_dir="/home/$user/.ssh"
        ssh_keys="$ssh_dir/authorized_keys"

        mkdir "$ssh_dir"
        cat "$file" > "$ssh_keys"
        chmod 700 "$ssh_dir"
        chmod 600 "$ssh_keys"
        chown -R "$user:$user" "$ssh_dir"


        # creates data dirs
        user_save=${SAVE_DIR}/${user}
        user_repo=${REPO_DIR}/${user}

        mkdir "$user_save" "$user_repo"
        chmod 0700 "$user_save" "$user_repo"
        chown "$user:$user" "$user_save" "$user_repo"

        echo "created user $user"
    done
}

create_users_loop() {
    while true; do
        create_users
        sleep 10
    done
}





# creates host key if it does not exist
if [ ! -f /etc/ssh/keys/ssh_host_ed25519_key ]; then
    echo "created host key"
    ssh-keygen -A
    mv /etc/ssh/ssh_host_*key* /etc/ssh/keys/
fi




mkdir -p "$REPO_DIR" "$SAVE_DIR"

create_users

# exec replaces the shell with the process sshd so that pid 1 is sshd and docker can stop it cleanly
# -D do not detach
# -e send error to standard error instead of /var/log, docker does not read /var log
exec /usr/sbin/sshd -D -e
