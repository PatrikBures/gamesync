#!/usr/bin/env bash

TESTDIR="./test"

USERS=(bob bertil bamse)
DEVICES=(test1 test2)
HOST_KEY_DIR="$TESTDIR/keys_host"
PRIVATE_KEY_DIR="$TESTDIR/keys_private"
PUBLIC_KEY_DIR="$TESTDIR/keys_public"
DATA="$TESTDIR/data"

CONTAINER_NAME="gamesync-test"

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

            echo "created key for $user for device $device"
        done
    done
fi

mkdir -p "$HOST_KEY_DIR"

docker container stop ${CONTAINER_NAME}
docker container rm ${CONTAINER_NAME}

docker build ./ -t gamesync:latest  || exit

docker run -d \
    --volume "$HOST_KEY_DIR":/etc/ssh/keys \
    --volume "$DATA":/data \
    -e "GAMESYNC_LOOP=10" \
    -e "GAMESYNC_UNRESTRICTED=true" \
    -p 127.0.0.1:2222:22 \
    --name ${CONTAINER_NAME} \
    gamesync:latest

sleep 0.2
docker logs ${CONTAINER_NAME}

echo "test server with:
docker exec -it ${CONTAINER_NAME} sh
"
echo "test ssh with:
ssh -i ./test/keys_private/bob_test1 -p 2222 gamesync@localhost
"
