package dbx

import (
	"context"
	"fmt"
	"gamesync/internal/server/dbm"
	"log/slog"
	"sync/atomic"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	primaryPool    *pgxpool.Pool
	replicaPools   []*pgxpool.Pool
	primaryQueries *dbm.Queries
	replicaQueries []*dbm.Queries
	readCounter    atomic.Int32
}

func NewDB(ctx context.Context, primaryConnStr string, replicaConnStrs []string) (*DB, error) {
	primaryPool, err := pgxpool.New(ctx, primaryConnStr)
	if err != nil {
		return nil, err
	}
	primaryQueries := dbm.New(primaryPool)

	replicaPools := make([]*pgxpool.Pool, 0, len(replicaConnStrs))
	replicaQueries := make([]*dbm.Queries, 0, len(replicaConnStrs))

	for _, cs := range replicaConnStrs {
		fmt.Println(len(replicaConnStrs))
		pool, err := pgxpool.New(ctx, cs)
		if err != nil {
			primaryPool.Close()
			for _, p := range replicaPools {
				p.Close()
			}
			return nil, err
		}
		replicaPools = append(replicaPools, pool)
		replicaQueries = append(replicaQueries, dbm.New(pool))
	}

	return &DB{
		primaryPool:    primaryPool,
		replicaPools:   replicaPools,
		primaryQueries: primaryQueries,
		replicaQueries: replicaQueries,
	}, nil
}

func (db *DB) ReadQuery() *dbm.Queries {
	// in case you have no replicas
	if len(db.replicaQueries) == 0 {
		return db.primaryQueries
	}

	counter := db.readCounter.Add(1)
	idx := (counter - 1) % int32(len(db.replicaQueries))

	return db.replicaQueries[idx]
}

func (db *DB) WriteQuery() *dbm.Queries {
	return db.primaryQueries
}

// returns queries which is used to run queries and a pgx.Tx
// which needs to be commited or rollbacked
func (db *DB) BeginTX(ctx context.Context) (*dbm.Queries, pgx.Tx, error) {
	tx, err := db.primaryPool.Begin(ctx)
	if err != nil {
		slog.Error("starting tx", "error", err)
		return nil, nil, err
	}
	qtx := db.WriteQuery().WithTx(tx)
	return qtx, tx, nil
}

func (db *DB) BeginReadTX(ctx context.Context) (*dbm.Queries, pgx.Tx, error) {
	opts := pgx.TxOptions{AccessMode: pgx.ReadOnly}
	if len(db.replicaQueries) == 0 {
		tx, err := db.primaryPool.BeginTx(ctx, opts)
		if err != nil {
			slog.Error("starting read tx", "errror", err)
			return nil, nil, err
		}
		qtx := db.primaryQueries.WithTx(tx)
		return qtx, tx, nil
	}

	counter := db.readCounter.Add(1)
	idx := (counter - 1) % int32(len(db.replicaQueries))
	tx, err := db.replicaPools[idx].BeginTx(ctx, opts)
	if err != nil {
		return nil, nil, err
	}
	qtx := dbm.New(tx)
	return qtx, tx, nil
}

// Closes write pool and all read pools
func (db *DB) Close() {
	db.primaryPool.Close()
	for _, p := range db.replicaPools {
		p.Close()
	}
}
