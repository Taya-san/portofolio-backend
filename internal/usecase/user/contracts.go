package user

import (
	"context"
)

type (
	User interface {
		Login(ctx context.Context, userItem UserItem) (string, error)
	}

	UserItem struct {
		Password string
	}
)
