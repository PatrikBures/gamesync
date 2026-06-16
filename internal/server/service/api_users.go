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

func (s *Service) GetUser(ctx context.Context, params api.GetUserParams) (api.GetUserRes, error) {
	// TODO: Seperate endpoint like GET /users/me to get info about the authenticated user

	user, err := s.db.ReadQuery().GetUser(ctx, params.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &api.GetUserNotFound{}, nil
		}
		return &api.GetUserInternalServerError{}, server.ErrDatabase
	}
	return &api.User{UserID: user.UserID, UserName: user.UserName, RoleID: user.RoleID}, nil
}

func (s *Service) GetUsers(ctx context.Context) (api.GetUsersRes, error) {
	users, err := s.db.ReadQuery().ListUsers(ctx)
	if err != nil {
		return &api.GetUsersInternalServerError{}, server.ErrDatabase
	}
	usersReturn := make(api.GetUsersOKApplicationJSON, 0, len(users))
	for _, user := range users {
		usersReturn = append(usersReturn, api.User{UserID: user.UserID, UserName: user.UserName, RoleID: user.RoleID})
	}

	return &usersReturn, nil
}

func (s *Service) PostUsers(ctx context.Context, req api.OptUserName) (result api.PostUsersRes, err error) {
	if !req.Set {
		return &api.PostUsersNotAcceptable{}, nil
	}
	user := dbm.InsertUserParams{
		UserName: req.Value.UserName,
		RoleID:   s.o.DefaultRoleID,
	}

	userID, token64, err := dbx.CreateUser(s.db, ctx, user)
	if err != nil {
		switch err {
		case server.ErrDatabase:
			return &api.PostUsersInternalServerError{}, err
		case server.ErrDuplicateKey:
			return &api.PostUsersConflict{}, nil
		case server.ErrToken:
			return &api.PostUsersInternalServerError{}, nil
		}
	}

	return &api.UserNewReturn{
		UserID:   userID,
		UserName: user.UserName,
		RoleID:   user.RoleID,
		Token:    token64,
	}, nil
}

func (s *Service) PutUserName(ctx context.Context, req api.OptUserName, params api.PutUserNameParams) (api.PutUserNameRes, error) {
	if !req.Set {
		return &api.PutUserNameNotAcceptable{}, nil
	}

	if err := s.db.WriteQuery().UpdateUserName(ctx, dbm.UpdateUserNameParams{
		UserID:   params.UserID,
		UserName: req.Value.UserName,
	}); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return &api.PutUserNameConflict{}, nil
			case pgerrcode.ForeignKeyViolation:
				return &api.PutUserNameNotFound{}, nil
			}
		}
		return &api.PutUserNameInternalServerError{}, server.ErrDatabase
	}

	return &api.PutUserNameOK{}, nil
}
