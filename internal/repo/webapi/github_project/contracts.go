package githubproject

import (
	"context"

	"github.com/Taya-san/portofolio-backend/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=.../usecase/mocks_webapi_githubproject_test.go -package=usecase_test

type (
	GithubProjectApi interface {
		Fetch(ctx context.Context) ([]entity.GithubProject, error)
	}
)
