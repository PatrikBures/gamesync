#!/usr/bin/env bash

TESTDIR="./test"

USERS=(bob bertil bamse)
PRIVATE_KEY_DIR="$TESTDIR/keys_private"
PUBLIC_KEY_DIR="$TESTDIR/keys_public"
HOST_KEY_DIR="$TESTDIR/keys_host"


if [ -d "$PRIVATE_KEY_DIR" ] || [ -d "$PUBLIC_KEY_DIR" ]; then
    mkdir -p "$PRIVATE_KEY_DIR"
    mkdir -p "$PUBLIC_KEY_DIR"

    for user in "${USERS[@]}"; do
        file="${PUBLIC_KEY_DIR}/${user}"
        ssh-keygen -C "${user}@test" -f "$file" -N "" > /dev/null
        mv "$file" "$PRIVATE_KEY_DIR"
        echo "created key for $user"
    done
fi

mkdir -p "$HOST_KEY_DIR"

docker container stop game-sync-test
docker container rm game-sync-test

docker build ./dockerfiles -t game-sync:latest 

docker run -d \
    --volume "$PUBLIC_KEY_DIR":/config/users \
    --volume "$HOST_KEY_DIR":/etc/ssh/keys \
    -p 127.0.0.1:2222:22 \
    --name game-sync-test \
    game-sync:latest

sleep 0.2
docker logs game-sync-test

echo "try ssh with:"
echo "ssh -i $PRIVATE_KEY_DIR/${USERS[0]} -p 2222 ${USERS[0]}@localhost"
