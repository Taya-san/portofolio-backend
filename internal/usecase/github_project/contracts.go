package githubproject

import (
	"context"

	"github.com/Taya-san/portofolio-backend/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../mocks_usecase_githubproject_test.go -package=usecase_test

type (
	GithubProject interface {
		RefreshMirror(ctx context.Context) error
		GetByProjectUrl(ctx context.Context, projectUrl string) (entity.GithubProject, error)
	}
)
