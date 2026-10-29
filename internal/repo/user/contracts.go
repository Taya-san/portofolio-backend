package user

import (
	"context"

	"github.com/Taya-san/portofolio-backend/internal/entity"
)

type (
	UserRepo interface {
		ListPermittedPasswordHashes(ctx context.Context) ([]entity.User, error)
	}
)
