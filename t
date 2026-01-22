#!/usr/bin/env bash

TESTDIR="./test"

USERS=(bob bertil bamse)
DEVICES=(test1 test2)
PRIVATE_KEY_DIR="$TESTDIR/keys_private"
PUBLIC_KEY_DIR="$TESTDIR/keys_public"
HOST_KEY_DIR="$TESTDIR/keys_host"
USER_IDS_DIR="$TESTDIR/user_ids"
RESTIC_PASSWORDS_DIR="$TESTDIR/restic_passwords"
DATA="$TESTDIR/data"


if ! [ -d "$PRIVATE_KEY_DIR" ] || ! [ -d "$PUBLIC_KEY_DIR" ]; then
    mkdir -p "$PRIVATE_KEY_DIR"
    mkdir -p "$PUBLIC_KEY_DIR"

    for user in "${USERS[@]}"; do
        for device in "${DEVICES[@]}"; do
            file="${PUBLIC_KEY_DIR}/${user}_${device}_tmp"
            ssh-keygen -C "${user}@${device}" -f "$file" -N "" > /dev/null
            mv "$file" "$PRIVATE_KEY_DIR/${user}_${device}"
            cat "${file}.pub" >> "${PUBLIC_KEY_DIR}/${user}"
            rm "${file}.pub"

            echo "created keyfor $user for device $device"
        done
    done
fi

mkdir -p "$HOST_KEY_DIR"

docker container stop gamesync-test
docker container rm gamesync-test

make build-state

docker build ./ -t gamesync:latest 

docker run -d \
    --volume "$USER_IDS_DIR":/config/user_ids \
    --volume "$PUBLIC_KEY_DIR":/config/users \
    --volume "$HOST_KEY_DIR":/etc/ssh/keys \
    --volume "$DATA":/data \
    -e "GAMESYNC_LOOP=10" \
    -p 127.0.0.1:2222:22 \
    --name gamesync-test \
    gamesync:latest

sleep 0.2
docker logs gamesync-test

echo "try ssh with:"
echo "ssh -i $PRIVATE_KEY_DIR/${USERS[0]}_${DEVICES[0]} -p 2222 ${USERS[0]}@localhost"
