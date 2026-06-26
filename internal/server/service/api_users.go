package service

import (
	"context"
	"errors"
	"gamesync/internal/dbx"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"
	"gamesync/internal/server/dbm"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Service) GetUser(ctx context.Context, params api.GetUserParams) (*api.User, error) {
	// TODO: Seperate endpoint like GET /users/me to get info about the authenticated user

	user, err := s.db.ReadQuery().GetUser(ctx, params.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, server.ErrUserNotFound
		}
		return nil, server.NewInternalError(err, "failed getting user", "userID", params.UserID)
	}
	return &api.User{UserID: user.UserID, UserName: user.UserName, RoleID: user.RoleID}, nil
}

func (s *Service) GetUsers(ctx context.Context) ([]api.User, error) {
	users, err := s.db.ReadQuery().ListUsers(ctx)
	if err != nil {
		return nil, server.NewInternalError(err, "failed getting users")
	}
	usersReturn := make([]api.User, 0, len(users))
	for _, user := range users {
		usersReturn = append(usersReturn, api.User{UserID: user.UserID, UserName: user.UserName, RoleID: user.RoleID})
	}

	return usersReturn, nil
}

func (s *Service) PostUsers(ctx context.Context, req *api.UserName) (result *api.UserNewReturn, err error) {
	user := dbm.InsertUserParams{
		UserName: req.UserName,
		RoleID:   s.o.DefaultRoleID,
	}

	userID, token64, err := dbx.CreateUser(s.db, ctx, user)
	if err != nil {
		return nil, err
	}

	return &api.UserNewReturn{
		UserID:   userID,
		UserName: user.UserName,
		RoleID:   user.RoleID,
		Token:    token64,
	}, nil
}

func (s *Service) PutUserName(ctx context.Context, req *api.UserName, params api.PutUserNameParams) error {
	if err := s.db.WriteQuery().UpdateUserName(ctx, dbm.UpdateUserNameParams{
		UserID:   params.UserID,
		UserName: req.UserName,
	}); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return server.ErrUserNameConflict
			case pgerrcode.ForeignKeyViolation:
				return server.ErrUserNotFound
			}
		}
		return server.NewInternalError(err, "failed updating user name", "userID", params.UserID, "newUserName", req.UserName)
	}

	return nil
}
