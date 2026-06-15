package service

import (
	"context"
	"errors"
	"gamesync/internal/dbx"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"
	"gamesync/internal/server/dbm"
	"gamesync/internal/server/permissions"
	"log/slog"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Service) GetUser(ctx context.Context, params api.GetUserParams) (api.GetUserRes, error) {
	// TODO: Seperate endpoint like GET /users/me to get info about the authenticated user

	currentUser, err := isUserSelf(ctx, params.UserID)
	if errors.Is(err, server.ErrContext) {
		return &api.GetUserInternalServerError{}, err
	} else if err == nil {
		return &api.User{UserID: currentUser.UserID, UserName: currentUser.UserName, RoleID: currentUser.RoleID}, nil
	}

	if err := CheckPerm(ctx, permissions.PermUserGet); err != nil {
		slog.Info("")
		return &api.GetUserUnauthorized{}, nil
	}
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
		RoleID: s.o.DefaultRoleID,
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
		UserID: userID,
		UserName: user.UserName,
		RoleID: user.RoleID,
		Token: token64,
	}, nil
}


func (s *Service) PutUserName(ctx context.Context, req api.OptUserName, params api.PutUserNameParams) (api.PutUserNameRes, error) {
	if !req.Set {
		return &api.PutUserNameNotAcceptable{}, nil
	}

	if err := isUserSelfAndPerm(ctx, params.UserID, permissions.PermUserNameUpdate); err != nil {
		if errors.Is(err, server.ErrContext) {
			return &api.PutUserNameInternalServerError{}, err
		} else if errors.Is(err, server.ErrNotAuthorized) {
			return &api.PutUserNameUnauthorized{}, nil
		}
	}

	if err := s.db.WriteQuery().UpdateUserName(ctx, dbm.UpdateUserNameParams{
		UserID: params.UserID,
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

