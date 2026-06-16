package middlewares

import (
	"gamesync/internal/server"
	"gamesync/internal/server/permissions"
	"log/slog"

	"github.com/ogen-go/ogen/middleware"
)

func UserAuthz() middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		targetUserID, ok := req.Params.Path("userID")
		if !ok {
			return next(req)
		}

		var userID int64
		switch i := targetUserID.(type) {
		case int64:
			userID = i
		default:
			slog.Error("user id is not int64", "path", req.Raw.URL.Path)
			return middleware.Response{}, server.ErrAuth
		}

		ctx := req.Context

		switch req.OperationName {
		case "PutUserName":
			if err := UserOrPerm(ctx, userID, permissions.PermUserNameUpdate); err != nil {
				return middleware.Response{}, server.ErrAuth
			}
		case "GetUser":
			if err := UserOrPerm(ctx, userID, permissions.PermUserGet); err != nil {
				return middleware.Response{}, server.ErrAuth
			}
		case "asdfasdf":
			slog.Info("change role")
		default:
			if _, err := UserSelf(ctx, userID); err != nil {
				return middleware.Response{}, server.ErrAuth
			}
		}
		return next(req)
	}
}
