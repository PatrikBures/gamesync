package middlewares

import (
	"gamesync/internal/server"
	"gamesync/internal/server/permissions"
	"gamesync/internal/server/service"
	"log/slog"

	"github.com/ogen-go/ogen/middleware"
)

var operationPerms = map[string]permissions.Perm{
	"GetHealth":          permissions.PermAllAllowed,
	"PostUsers":          permissions.PermAllAllowed,
	"GetPerms":           permissions.PermAllAllowed,


	"GetUser":            permissions.PermUserGetOwn,
	"GetUsers":           permissions.PermUsersList,

	"PutUserName":        permissions.PermUserNameUpdateOwn,

	"GetRoles":           permissions.PermRolesGet,
	"PostRoles":          permissions.PermRolesMod,

	"PutRoleName":        permissions.PermRolesMod,
	"GetRolePerms":       permissions.PermRolesGet,
	"PatchRolePerms":     permissions.PermRolesMod,
	"PutRolePerms":       permissions.PermRolesMod,
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
		if err := service.CheckPerm(ctx, perm); err != nil {
			return middleware.Response{}, err
		}
		return next(req)
	}
}

