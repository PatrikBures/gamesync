ARG ALPINE_VERSION="3.23"

FROM golang:1.26-alpine AS build
WORKDIR /src

ENV GOCACHE=/cache/go
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
ENV CGO_ENABLED="0"
RUN --mount=type=cache,target="/cache/go" go build -ldflags="-s -w" -trimpath -o /bin/gamesync-admin   ./cmd/gamesync-admin
RUN --mount=type=cache,target="/cache/go" go build -ldflags="-s -w" -trimpath -o /bin/gamesync-auth    ./cmd/gamesync-auth
RUN --mount=type=cache,target="/cache/go" go build -ldflags="-s -w" -trimpath -o /bin/gamesync-wrapper ./cmd/gamesync-wrapper


FROM alpine:${ALPINE_VERSION}

RUN apk add --no-cache \
    rsync \
    openssh-server \
    tzdata \
    bash

COPY --from=build --chmod=555 /bin/gamesync-admin     /usr/local/bin/gamesync-admin
COPY --from=build --chmod=555 /bin/gamesync-auth      /usr/local/bin/gamesync-auth
COPY --from=build --chmod=555 /bin/gamesync-wrapper   /usr/local/bin/gamesync-wrapper

COPY server/sshd.conf /etc/ssh/sshd_config.d/50-game-sync.conf
COPY --chmod=500 server/entrypoint.sh /entrypoint.sh

EXPOSE 22

CMD ["/entrypoint.sh"]
