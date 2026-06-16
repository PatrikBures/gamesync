package initServer

import (
	"context"
	"errors"
	"gamesync/internal/dbx"
	"gamesync/internal/server/dbm"

	"github.com/jackc/pgx/v5"
)

func CreateAdmin(db *dbx.DB) (string, error) {
	user := dbm.InsertUserParams{
		UserName: "admin",
		RoleID:   1,
	}

	ctx := context.Background()

	curUser, err := db.ReadQuery().GetUserWithName(ctx, "admin")
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", err
	} else if curUser.RoleID > 0 {
		return "", nil
	}

	_, token, err := dbx.CreateUser(db, ctx, user)
	if err != nil {
		return "", err
	}
	return token, nil
}
