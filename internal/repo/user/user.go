package user

import (
	"context"

	"github.com/Taya-san/portofolio-backend/config"
	"github.com/Taya-san/portofolio-backend/internal/entity"
)

type Repo struct {
	*config.Config
}

func New(cfg *config.Config) UserRepo {
	return newTraced(&Repo{cfg})
}

func (r *Repo) ListPermittedPasswordHashes(ctx context.Context) ([]entity.User, error) {
	var users []entity.User

	for _, password := range r.User.Password {
		users = append(users, entity.User{
			Password: password,
		})
	}

	return users, nil
}
