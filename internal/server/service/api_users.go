package service

import (
	"context"
	"errors"
	"gamesync/internal/dbx"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"
	"gamesync/internal/server/dbm"
	"gamesync/internal/server/permissions"

	"github.com/jackc/pgx/v5"
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
		return &api.GetUserUnauthorized{}, err
	}
	user, err := s.conn.Queries.GetUser(ctx, params.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &api.GetUserNotFound{}, server.ErrNotFound
		}
		return &api.GetUserInternalServerError{}, server.ErrDatabase
	}
	return &api.User{UserID: user.UserID, UserName: user.UserName, RoleID: user.RoleID}, nil
}

func (s *Service) GetUsers(ctx context.Context) (api.GetUsersRes, error) {
	users, err := s.conn.Queries.ListUsers(ctx)
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
		return &api.PostUsersNotAcceptable{}, server.ErrMissingBody
	}
	user := dbm.InsertUserParams{
		UserName: req.Value.UserName,
		RoleID: s.o.DefaultRoleID,
	}

	userID, token64, err := dbx.CreateUser(s.conn, ctx, user)
	if err != nil {
		switch err {
		case server.ErrDatabase:
			return &api.PostUsersInternalServerError{}, err
		case server.ErrDuplicateKey:
			return &api.PostUsersConflict{}, err
		case server.ErrToken:
			return &api.PostUsersInternalServerError{}, err
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
		return &api.PutUserNameNotAcceptable{}, server.ErrMissingBody
	}

	// if err := isUserSelfAndPerm(ctx, params.UserID, permissions.PermUserNameUpdate); err != nil {
	// 	if errors.Is(err, server.ErrContext) {
	// 		return &api.PutUserNameInternalServerError{}, err
	// 	} else if errors.Is(err, server.ErrNotAuthorized) {
	// 		return &api.PutUserNameUnauthorized{}, err
	// 	}
	// }
	//
	// if _, err := s.q.User.WithContext(ctx).
	// 	Where(s.q.User.UserID.Eq(params.UserID)).
	// 	Update(s.q.User.UserName, req.Value.UserName);
	// 	err != nil {
	//
	// 	if errors.Is(err, gorm.ErrDuplicatedKey) {
	// 		return &api.PutUserNameConflict{}, server.ErrDuplicateKey
	// 	}
	// 	return &api.PutUserNameInternalServerError{}, server.ErrDatabase
	// }
	return &api.PutUserNameOK{}, nil
}

