#!/bin/bash

DATA_DIR=/data
SAVES_DIR=$DATA_DIR/saves
REPOS_DIR=$DATA_DIR/repos
CONFIG_DIR=/config
USERS_DIR=$CONFIG_DIR/users
USER_IDS_DIR=$CONFIG_DIR/user_ids
RESTIC_PASSWORDS_DIR=$CONFIG_DIR/restic_passwords

GROUP=client-users

random_password() {
    hexdump -n 64 -ve '/1 "%02x"' /dev/urandom
}

create_user() {
    user=$1
    id=$2

    if id -u "$user" &> /dev/null; then
        return
    fi

    if [[ -n "$id" ]]; then
        adduser "$user" -D -s /bin/sh -u "$id" > /dev/null
    else
        adduser "$user" -D -s /bin/sh > /dev/null
    fi

    adduser "$user" "$GROUP" > /dev/null
    passwd -u "$user" &> /dev/null


    # adds ssh key to authorized_keys
    ssh_dir="/home/$user/.ssh"
    ssh_keys="$ssh_dir/authorized_keys"

    mkdir "$ssh_dir"
    cat "$USERS_DIR/$user" > "$ssh_keys"
    chmod 700 "$ssh_dir"
    chmod 600 "$ssh_keys"
    chown -R "$user:$user" "$ssh_dir"


    # creates data dirs
    user_save=${SAVES_DIR}/${user}
    user_repo=${REPOS_DIR}/${user}

    mkdir -p "$user_save" "$user_repo"
    chmod -R 0700 "$user_save" "$user_repo"
    chown -R "$user:$user" "$user_save" "$user_repo"

    echo "created user $user $id"

    RESTIC_PASSWORD_FILE="$RESTIC_PASSWORDS_DIR/${user}"

    if ! [ -f "$RESTIC_PASSWORD_FILE" ]; then 
        echo "created restic password for $user at ${RESTIC_PASSWORD_FILE}"
        random_password > "$RESTIC_PASSWORD_FILE"
    fi

    chown "$user:$user" "$RESTIC_PASSWORD_FILE"
    chmod 400 "$RESTIC_PASSWORD_FILE"

    id -u "$user" > "$USER_IDS_DIR/$user"
}

create_users() {
    users_with_ids=()
    users_with_ids_id=()
    users=()

    for file in "$USERS_DIR"/*; do
        user="$(basename "${file}")"

        id_file="$USER_IDS_DIR/$user"

        if [[ -f "$id_file" ]]; then
            id=$(cat "$id_file")

            if [[ "$id" -lt 1000 ]]; then
                echo "user $user has invalid id: $id"
                exit 1
            fi

            users_with_ids_id+=("$id")
            users_with_ids+=("$user")
        else
            users+=("$user")
        fi
    done

    for i in "${!users_with_ids[@]}"; do
        user="${users_with_ids[$i]}"
        id="${users_with_ids_id[$i]}"

        create_user "$user" "$id"
    done


    # remove owners of removed users
    for file in "$SAVES_DIR"/*; do
        user=$(basename "$file")

        if [[ " ${users_with_ids[*]} ${users[*]} " =~ [[:space:]]${user}[[:space:]] ]]; then
            continue
        fi
        
        s="$SAVES_DIR/$user"
        r="$REPOS_DIR/$user"

        dir_owner=$(stat -c '%U' "$s")
        
        if [[ "$dir_owner" == "nobody" ]]; then
            continue
        fi

        chown -R nobody:nobody "$s" "$r"
        echo "made $user dirs owned by nobody"
    done

    for user in "${users[@]}"; do
        create_user "$user"
    done
}

create_users_loop() {
    while true; do
        sleep "$GAMESYNC_LOOP"
        create_users 
    done
}


if [[ "$GAMESYNC_UNRESTRICTED" == "true" ]]; then
    echo "export GAMESYNC_UNRESTRICTED=true" >> /etc/gamesync.env
    chmod 664 /etc/gamesync.env
fi



# creates host key if it does not exist
if [ ! -f /etc/ssh/keys/ssh_host_ed25519_key ]; then
    echo "created host key"
    ssh-keygen -A
    mv /etc/ssh/ssh_host_*key* /etc/ssh/keys/
fi



addgroup -S "$GROUP"

mkdir -p "$REPOS_DIR" "$SAVES_DIR" "$USER_IDS_DIR" "$RESTIC_PASSWORDS_DIR"

chmod -R 0550 "$USER_IDS_DIR"
chmod 0770 "$USER_IDS_DIR"

create_users

if [[ $GAMESYNC_LOOP -gt 0 ]]; then
    create_users_loop &
fi

# exec replaces the shell with the process sshd so that pid 1 is sshd and docker can stop it cleanly
# -D do not detach
# -e send error to standard error instead of /var/log, docker does not read /var log
exec /usr/sbin/sshd -D -e
