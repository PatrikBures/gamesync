FROM golang:1.26-alpine AS builder

WORKDIR /src

ADD https://github.com/pressly/goose.git .

RUN go build -tags='no_clickhouse no_libsql no_mssql no_mysql no_vertica no_ydb no_postgres' -o goose ./cmd/goose



FROM busybox:1

COPY --from=builder --chmod=775 /src/goose /usr/local/bin/goose

COPY migrations/sqlite/ /migrations

ENV GOOSE_DRIVER=sqlite3
ENV GOOSE_DBSTRING=/db/gamesync.sqlite
ENV GOOSE_MIGRATION_DIR=/migrations
