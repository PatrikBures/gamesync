package dbx

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"gamesync/internal/model"
	"gamesync/internal/query"
	"gamesync/internal/server"
	"log/slog"

	"gorm.io/gorm"
)

// error can be: ErrDatabase, ErrToken, ErrDuplicateKey or nil
func CreateUser(tx *query.QueryTx, ctx context.Context, user *model.User) (token64 string, err error) {
	defer func() {
		if recover() != nil || err != nil {
			if e := tx.Rollback(); e != nil {
				slog.Error("failed rollback", "error", e)
			}
		}
	}()

	if err = tx.User.WithContext(ctx).Create(user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return "", server.ErrDuplicateKey
		} else if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return "", server.ErrDatabase
		}
		return "", server.ErrDatabase
	}

	token, err := generateToken()
	if err != nil {
		return "", server.ErrToken
	}
	token64 = base64.URLEncoding.EncodeToString(token)
	tokenHash := sha256.Sum256(token)

	t := model.Token{
		UserID: user.UserID,
		TokenHash: tokenHash[:],
	}
	if err = tx.Token.WithContext(ctx).Create(&t); err != nil {
		return "", server.ErrDatabase
	}

	if err = tx.Commit(); err != nil {
		return "", server.ErrDatabase
	}

	return token64, nil
}


func generateToken() ([]byte, error) {
	b := make([]byte, 33)
	_, err := rand.Read(b)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
