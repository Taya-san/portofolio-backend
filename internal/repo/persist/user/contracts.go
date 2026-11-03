package user

import (
	"context"

	"github.com/Taya-san/portofolio-backend/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../../../usecase/mocks_repo_user_test.go -package=usecase_test

type (
	UserRepo interface {
		ListPermittedPassword(ctx context.Context) ([]entity.User, error)
	}
)
