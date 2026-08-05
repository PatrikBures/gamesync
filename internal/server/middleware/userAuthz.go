package middlewares

import (
	"go.pabu.dev/gamesync/internal/server"
	"go.pabu.dev/gamesync/internal/server/permissions"
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
			return middleware.Response{}, server.NewInternalError(nil, "userID is somehow not int64", "path", req.Raw.URL)
		}

		ctx := req.Context

		switch req.OperationName {
		case "PutUserName":
			if err := UserOrPerm(ctx, userID, permissions.PermUserNameUpdate); err != nil {
				return middleware.Response{}, server.ErrNotAuthorized
			}
		case "GetUser":
			if err := UserOrPerm(ctx, userID, permissions.PermUserGet); err != nil {
				return middleware.Response{}, server.ErrNotAuthorized
			}
		case "asdfasdf":
			slog.Info("change role")
		default:
			if _, err := UserSelf(ctx, userID); err != nil {
				return middleware.Response{}, server.ErrNotAuthorized
			}
		}
		return next(req)
	}
}
