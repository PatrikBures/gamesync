package middlewares

import (
	"context"
	"errors"
	"gamesync/internal/server"
	"gamesync/internal/server/dbm"
	"gamesync/internal/server/permissions"
	"gamesync/internal/server/service"
	"log/slog"
	"slices"

	"github.com/ogen-go/ogen/middleware"
)

var operationPerms = map[string]permissions.Perm{
	"GetHealth": permissions.PermAllAllowed,
	"PostUsers": permissions.PermAllAllowed,
	"GetPerms":  permissions.PermAllAllowed,

	"GetUser":  permissions.PermUserGetOwn,
	"GetUsers": permissions.PermUsersList,

	"PutUserName": permissions.PermUserNameUpdateOwn,

	"GetRoles":  permissions.PermRolesGet,
	"PostRoles": permissions.PermRolesMod,

	"PutRoleName":    permissions.PermRolesMod,
	"GetRolePerms":   permissions.PermRolesGet,
	"PatchRolePerms": permissions.PermRolesMod,
	"PutRolePerms":   permissions.PermRolesMod,

	"GetUserRepos": permissions.PermSync,
	"PutUserRepo":  permissions.PermSync,

	"GetUserRepoBranches": permissions.PermSync,
	"PutUserRepoBranch":   permissions.PermSync,
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
		slog.Error("failed mapping role perms context")
		return server.ErrContext
	}
	if !slices.Contains(perms, perm) {
		return server.ErrNotAuthorized
	}
	return nil
}

func UserSelf(ctx context.Context, userID int64) (dbm.User, error) {
	currentUser, ok := ctx.Value(service.CkUser).(dbm.User)
	if !ok {
		slog.Error("failed mapping user context", "UserID", userID)
		return dbm.User{}, server.ErrContext
	}
	if userID != currentUser.UserID {
		return currentUser, server.ErrNotAuthorized
	}
	return currentUser, nil
}
// func UserSelf(u)

// checks if the user is trying to access itself,
// if it is it will return nil
//
// if it is trying to access a user id which is not itself,
// it will check if it has the perm for that
//
// returns either ErrContext, ErrNotAuthorized or nil
func UserOrPerm(ctx context.Context, userID int64, perm permissions.Perm) error {
	_, errUser := UserSelf(ctx, userID)
	if errors.Is(errUser, server.ErrContext) {
		return errUser
	} else if errUser == nil {
		return nil
	}

	errPerm := CheckPerm(ctx, perm)
	if errors.Is(errPerm, server.ErrContext) {
		return errPerm
	} else if errPerm != nil {
		return server.ErrNotAuthorized
	}

	return nil
}

