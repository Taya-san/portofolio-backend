package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Taya-san/portofolio-backend/internal/entity"
	"github.com/Taya-san/portofolio-backend/internal/usecase/user"
	"github.com/Taya-san/portofolio-backend/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newUserUseCase(t *testing.T) (user.User, *MockUserRepo) {
	t.Helper()

	ctrl := gomock.NewController(t)

	repo := NewMockUserRepo(ctrl)
	jwtManager := jwt.New("test-secret", time.Hour)
	useCase := user.New(repo, jwtManager)

	return useCase, repo
}

func TestLogin(t *testing.T) {
	t.Parallel()

	t.Run("login success", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)

		storedUser := entity.User{
			Password: "098",
		}

		var storedUsers []entity.User

		storedUsers = append(storedUsers, storedUser)
		repo.EXPECT().ListPermittedPassword(gomock.Any()).Return(storedUsers, nil)

		token, err := uc.Login(context.Background(), user.UserItem{Password: "098"})

		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("login wrong password", func(t *testing.T) {
		t.Parallel()

		uc, repo := newUserUseCase(t)

		storedUser := entity.User{
			Password: "098",
		}

		var storedUsers []entity.User

		storedUsers = append(storedUsers, storedUser)
		repo.EXPECT().ListPermittedPassword(gomock.Any()).Return(storedUsers, nil)

		token, err := uc.Login(context.Background(), user.UserItem{Password: "wrongpassword"})

		require.ErrorIs(t, err, entity.ErrInvalidCredentials)
		assert.Empty(t, token)
	})
}
