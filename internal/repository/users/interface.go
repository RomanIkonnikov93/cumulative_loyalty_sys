package users

import (
	"context"
)

type Users interface {
	AddUserAuthData(ctx context.Context, login, pass, token string) error
	GetUserAuthData(ctx context.Context, login, pass string) (string, error)
}
