ARG ALPINE_VERSION="3.23"

FROM golang:1.26-alpine AS build
WORKDIR /src

ENV GOCACHE=/cache/go
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
ENV CGO_ENABLED="0"
RUN --mount=type=cache,target="/cache/go" go build -ldflags="-s -w" -trimpath -o /bin/gamesync-server ./cmd/gamesync-server


FROM alpine:${ALPINE_VERSION}

RUN apk add --no-cache \
    rsync \
    openssh-server

COPY --from=build --chmod=555 /bin/gamesync-server /usr/local/bin/gamesync-server

RUN ln -s /usr/local/bin/gamesync-server /usr/local/bin/gamesync-admin && \
    ln -s /usr/local/bin/gamesync-server /usr/local/bin/gamesync-auth && \
    ln -s /usr/local/bin/gamesync-server /usr/local/bin/gamesync-wrapper

COPY server/sshd.conf /etc/ssh/sshd_config.d/50-game-sync.conf
COPY --chmod=500 server/entrypoint.sh /entrypoint.sh

EXPOSE 22

CMD ["/entrypoint.sh"]
