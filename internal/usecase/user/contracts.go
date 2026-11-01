package user

import (
	"context"
)

//go:generate mockgen -source=contracts.go -destination=../mocks_usecase_user_test.go -package=usecase_test

type (
	User interface {
		Login(ctx context.Context, userItem UserItem) (string, error)
	}

	UserItem struct {
		Password string
	}
)
