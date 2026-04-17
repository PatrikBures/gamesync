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

EXPOSE 22

CMD ["/usr/local/bin/gamesync-server"]
