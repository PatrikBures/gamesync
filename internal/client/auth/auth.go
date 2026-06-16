package auth

import (
	"context"
	api "gamesync/internal/ogen"
)

type auth struct {
	token string
}

func NewAuth(token string) *auth {
	return &auth{token: token}
}
func (a *auth) BearerAuth(ctx context.Context, operationName api.OperationName) (api.BearerAuth, error) {
	return api.BearerAuth{Token: a.token}, nil
}
