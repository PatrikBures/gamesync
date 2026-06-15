package dbx

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"gamesync/internal/server"
	"gamesync/internal/server/dbm"
	"log/slog"
)

// error can be: ErrDatabase, ErrToken, ErrDuplicateKey or nil
func CreateUser(db *DB, ctx context.Context, user dbm.InsertUserParams) (userID int64, token64 string, err error) {
	qtx, tx, err := db.BeginTX(ctx)
	if err != nil {
		slog.Error("starting tx", "error", err)
		return
	}

	defer func() {
		if recover() != nil || err != nil {
			if e := tx.Rollback(ctx); e != nil {
				slog.Error("failed rollback", "error", e)
			}
		}
	}()

	userID, err = qtx.InsertUser(ctx, user)
	if err != nil {
		slog.Error("inserting user", "user", user, "error", err)
		return 0, "", server.ErrDatabase
	}

	token, err := generateToken()
	if err != nil {
		slog.Error("generating token", "error", err)
		return 0, "", server.ErrToken
	}
	token64 = base64.URLEncoding.EncodeToString(token)
	tokenHash := sha256.Sum256(token)

	t := dbm.InsertTokenParams{UserID: userID, TokenHash: tokenHash[:]}
	_, err = qtx.InsertToken(ctx, t)
	if err != nil {
		slog.Error("inserting token", "token", t, "error", err)
		return 0, "", server.ErrDatabase
	}

	if err = tx.Commit(ctx); err != nil {
		slog.Error("commiting tx", "error", err)
		return 0, "", server.ErrDatabase
	}

	return
}


func generateToken() ([]byte, error) {
	b := make([]byte, 33)
	_, err := rand.Read(b)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
