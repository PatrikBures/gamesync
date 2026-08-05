package dbx

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"go.pabu.dev/gamesync/internal/server"
	"go.pabu.dev/gamesync/internal/server/dbm"
	"log/slog"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// error can be: InternalError, ErrUserNameConflict or nil
func CreateUser(db *DB, ctx context.Context, user dbm.InsertUserParams) (userID int64, token64 string, err error) {
	qtx, tx, err := db.BeginTX(ctx)
	if err != nil {
		return 0, "", server.NewInternalError(err, "failed starting tx")
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
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return 0, "", server.ErrUserNameConflict
			}
		}

		return 0, "", server.NewInternalError(err, "failed inserting new user")
	}

	token, err := generateToken()
	if err != nil {
		return 0, "", server.NewInternalError(err, "failed generating new token", "userID", userID)
	}
	token64 = base64.URLEncoding.EncodeToString(token)
	tokenHash := sha256.Sum256(token)

	t := dbm.InsertTokenParams{UserID: userID, TokenHash: tokenHash[:]}
	_, err = qtx.InsertToken(ctx, t)
	if err != nil {
		return 0, "", server.NewInternalError(err, "failed inserting token", "userID", userID)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, "", server.NewInternalError(err, "committing tx when creating user", "userID", userID)
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
