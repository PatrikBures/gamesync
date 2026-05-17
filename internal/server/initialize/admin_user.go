package initServer

import (
	"context"
	"gamesync/internal/dbx"
	"gamesync/internal/model"
	"gamesync/internal/query"
)

func CreateAdmin(q *query.Query) (string, error) {
	user := &model.User{
		UserID: 1,
		UserName: "admin",
		RoleID: 1,
	}

	ctx := context.Background()

	existingUserCount, err := q.User.WithContext(ctx).Where(q.User.UserID.Eq(user.UserID)).Count()
	if err != nil {
		return "", err
	}
	if existingUserCount == 1 {
		return "", nil
	}

	tx := q.Begin()

	token, err := dbx.CreateUser(tx, ctx, user)
	if err != nil {
		return "", err
	}
	return token, nil
}
