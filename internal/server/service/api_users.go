package service

import (
	"context"
	"gamesync/internal/dbx"
	"gamesync/internal/model"
	api "gamesync/internal/ogen"
	"gamesync/internal/server"
	"gamesync/internal/server/permissions"
	"log/slog"
)

func (s *Service) GetUser(ctx context.Context, params api.GetUserParams) (api.GetUserRes, error) {
	// TODO: Seperate endpoint like GET /users/me to get info about the authenticated user
	currentUser, ok := ctx.Value(ckUser).(*model.User)
	if !ok {
		slog.Error("mapping user context", "UserID", params.UserID)
		return &api.GetUserInternalServerError{}, server.ErrContext
	}
	if params.UserID == currentUser.UserID {
		if err := CheckPerm(ctx, permissions.PermUserGetOwn); err != nil {
			return &api.GetUserUnauthorized{}, err
		}
		return &api.User{UserID: currentUser.UserID, UserName: currentUser.UserName, RoleID: currentUser.RoleID}, nil
	}

	if err := CheckPerm(ctx, permissions.PermUserGet); err != nil {
		return &api.GetUserUnauthorized{}, err
	}

	user, err := s.q.User.WithContext(ctx).Where(s.q.User.UserID.Eq(params.UserID)).First()
	if err != nil {
		return &api.GetUserInternalServerError{}, server.ErrDatabase
	}
	return &api.User{UserID: user.UserID, UserName: user.UserName, RoleID: user.RoleID}, nil
}

func (s *Service) GetUsers(ctx context.Context) (api.GetUsersRes, error) {
	users, err := s.q.User.WithContext(ctx).Find()
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
	user := &model.User{
		UserName: req.Value.UserName,
		RoleID: s.o.DefaultRoleID,
	}
	tx := s.q.Begin()

	token64, err := dbx.CreateUser(tx, ctx, user)
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
		UserID: user.UserID,
		UserName: user.UserName,
		RoleID: user.RoleID,
		Token: token64,
	}, nil
}


