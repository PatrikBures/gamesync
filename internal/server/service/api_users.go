package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"gamesync/internal/model"
	api "gamesync/internal/ogen"

	"gorm.io/gorm"
)

func (s *Service) GetUser(ctx context.Context, params api.GetUserParams) (api.GetUserRes, error) {
	user, err := s.q.User.WithContext(ctx).Where(s.q.User.UserID.Eq(params.UserID)).First()
	if err != nil {
		return &api.GetUserInternalServerError{}, ErrDatabase
	}
	return &api.User{UserID: user.UserID, UserName: user.UserName, RoleID: user.RoleID}, nil
}

func (s *Service) GetUsers(ctx context.Context) (api.GetUsersRes, error) {
	users, err := s.q.User.WithContext(ctx).Find()
	if err != nil {
		return &api.GetUsersInternalServerError{}, ErrDatabase
	}
	usersReturn := make(api.GetUsersOKApplicationJSON, 0, len(users))
	for _, user := range users {
		usersReturn = append(usersReturn, api.User{UserID: user.UserID, UserName: user.UserName, RoleID: user.RoleID})
	}

	return &usersReturn, nil
}

func (s *Service) PostUsers(ctx context.Context, req api.OptUserName) (result api.PostUsersRes, err error) {
	if !req.Set {
		return &api.PostUsersNotAcceptable{}, ErrMissingBody
	}
	user := model.User{
		UserName: req.Value.UserName,
		RoleID: s.o.DefaultRoleID,
	}
	tx := s.q.Begin()
	defer func() {
		if recover() != nil || err != nil {
			err = tx.Rollback()
		}
	}()

	if err = tx.User.WithContext(ctx).Create(&user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &api.PostUsersConflict{}, ErrDuplicateKey
		} else if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return &api.PostUsersNotAcceptable{}, fmt.Errorf("invalid role")
		}
		return &api.PostUsersInternalServerError{}, ErrDatabase
	}

	token, err := generateToken()
	if err != nil {
		return &api.PostUsersInternalServerError{}, ErrToken
	}
	token64 := base64.URLEncoding.EncodeToString(token)
	tokenHash := sha256.Sum256(token)

	t := model.Token{
		UserID: user.UserID,
		TokenHash: tokenHash[:],
	}
	if err = tx.Token.WithContext(ctx).Create(&t); err != nil {
		return &api.PostUsersInternalServerError{}, ErrDatabase
	}

	if err = tx.Commit(); err != nil {
		return &api.PostUsersInternalServerError{}, ErrDatabase
	}

	return &api.UserNewReturn{
		UserID: user.UserID,
		UserName: user.UserName,
		RoleID: user.RoleID,
		Token: token64,
	}, nil
}
