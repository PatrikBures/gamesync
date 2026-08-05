package middlewares

import (
	"context"
	"errors"
	"go.pabu.dev/gamesync/internal/server"
	"go.pabu.dev/gamesync/internal/server/dbm"
	"go.pabu.dev/gamesync/internal/server/permissions"
	"go.pabu.dev/gamesync/internal/server/service"
	"log/slog"
	"slices"

	"github.com/ogen-go/ogen/middleware"
)

var operationPerms = map[string]permissions.Perm{
	"GetHealth": permissions.PermAllAllowed,
	"PostUsers": permissions.PermAllAllowed,
	"GetPerms":  permissions.PermAllAllowed,
	"GetMe":     permissions.PermAllAllowed,

	"GetUser":  permissions.PermUserGetOwn,
	"GetUsers": permissions.PermUsersList,

	"PutUserName": permissions.PermUserNameUpdateOwn,

	"GetRoles":  permissions.PermRolesGet,
	"PostRoles": permissions.PermRolesMod,

	"PutRoleName":    permissions.PermRolesMod,
	"GetRolePerms":   permissions.PermRolesGet,
	"PatchRolePerms": permissions.PermRolesMod,
	"PutRolePerms":   permissions.PermRolesMod,

	"GetRepos": permissions.PermSync,
	"PutRepo":  permissions.PermSync,

	"GetBranches": permissions.PermSync,
	"PutBranch":   permissions.PermSync,
	"GetBranchHead":       permissions.PermSync,

	"PostSnapshot": permissions.PermSync,
	"GetSnapshot":         permissions.PermSync,

	"PutChunk":            permissions.PermSync,
	"GetChunk":            permissions.PermSync,
}

func AuthzMiddleware() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		perm, ok := operationPerms[req.OperationName]
		if !ok {
			slog.Error("no permission mapping for operation", "operation", req.OperationName)
			return middleware.Response{}, server.ErrNotAuthorized
		}

		if perm == permissions.PermAllAllowed {
			return next(req)
		}

		ctx := req.Context
		if err := CheckPerm(ctx, perm); err != nil {
			slog.Warn("checking permission", "perm", perm, "err", err)
			return middleware.Response{}, err
		}
		return next(req)
	}
}



// returns either ErrContext, ErrNotAuthorized or nil
func CheckPerm(ctx context.Context, perm permissions.Perm) error {
	perms, ok := ctx.Value(service.CkRolePerms).(permissions.Perms)
	if !ok {
		return server.NewInternalError(nil, "failed mapping role perms context")
	}
	if !slices.Contains(perms, perm) {
		return server.ErrNotAuthorized
	}
	return nil
}

func UserSelf(ctx context.Context, userID int64) (dbm.User, error) {
	currentUser, ok := ctx.Value(service.CkUser).(dbm.User)
	if !ok {
		return dbm.User{}, server.NewInternalError(nil, "failed mapping role perms context", "userID", userID)
	}
	if userID != currentUser.UserID {
		return currentUser, server.ErrNotAuthorized
	}
	return currentUser, nil
}

// checks if the user is trying to access itself,
// if it is it will return nil
//
// if it is trying to access a user id which is not itself,
// it will check if it has the perm for that
//
// returns either InternalError, ErrNotAuthorized or nil
func UserOrPerm(ctx context.Context, userID int64, perm permissions.Perm) error {
	_, errUser := UserSelf(ctx, userID)
	if errUser == nil {
		return nil
	} else if !errors.Is(errUser, server.ErrNotAuthorized) {
		return errUser
	}

	return CheckPerm(ctx, perm) 
}

