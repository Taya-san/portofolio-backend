package user

import (
	"context"
	"fmt"

	"github.com/Taya-san/portofolio-backend/internal/repo/user"
	"github.com/Taya-san/portofolio-backend/pkg/jwt"
)

type UseCase struct {
	repo user.UserRepo
	jwt  *jwt.Manager
}

func New(r user.UserRepo, j *jwt.Manager) User {
	return newTraced(&UseCase{
		repo: r,
		jwt:  j,
	})
}

func (uc *UseCase) Login(ctx context.Context, userItem UserItem) (string, error) {
	users, err := uc.repo.ListPermittedPasswordHashes(ctx)
	if err != nil {
		return "", fmt.Errorf("UserUseCase - Login - uc.repo.ListPermittedPasswordHashes: %w", err)
	}

	for _, user := range users {
		if user.Password != userItem.Password {
			continue
		}

		token, err := uc.jwt.GenerateToken()
		if err != nil {
			return "", fmt.Errorf("UserUseCase - Login - uc.jwt.GenerateToken: %w", err)
		}

		return token, nil
	}

	return "", fmt.Errorf("UserUseCase - Login - User Not Found: %w", err)
}
