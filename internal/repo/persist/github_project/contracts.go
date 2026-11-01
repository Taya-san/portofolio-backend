package githubproject

import (
	"context"

	"github.com/Taya-san/portofolio-backend/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=.../usecase/mocks_repo_githubproject_test.go -package=usecase_test

type (
	GithubProjectRepo interface {
		UpStore(ctx context.Context, project *entity.GithubProject) error
		GetByProjectUrl(ctx context.Context, projectUrl string) (entity.GithubProject, error)
		List(ctx context.Context) ([]entity.GithubProject, int, error)
		Delete(ctx context.Context, projectID string) error
	}
)
