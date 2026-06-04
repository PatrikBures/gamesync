package dbx

import (
	"context"
	"gamesync/internal/server/dbm"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBconn struct {
	Queries *dbm.Queries
	Pool *pgxpool.Pool
}
func ConnectDb(dsn string) (DBconn, error) {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)

	return DBconn{Queries: dbm.New(pool), Pool: pool}, err
}
