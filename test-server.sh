#!/usr/bin/env bash

TESTDIR="./test"

HOST_KEY_DIR="$TESTDIR/keys_host"
DATA="$TESTDIR/data"

CONTAINER_NAME="gamesync-test"


mkdir -p "$HOST_KEY_DIR"

docker container stop ${CONTAINER_NAME}
docker container rm ${CONTAINER_NAME}

make build-admin
make build-state

docker build ./ -t gamesync:latest 

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
docker exec -it ${CONTAINER_NAME} bash"
