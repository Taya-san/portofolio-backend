package project

import (
	"context"

	"github.com/Taya-san/portofolio-backend/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_project_test.go -package=usecase_test

type (
	ProjectRepo interface {
		Store(ctx context.Context, project *entity.Project) error
		GetByID(ctx context.Context, projectID string) (entity.Project, error)
		List(ctx context.Context, filter ProjectFilter) ([]entity.Project, int, error)
		Update(ctx context.Context, project *entity.Project) error
		Delete(ctx context.Context, projectID string) error
	}

	ProjectFilter struct {
		All bool
	}
)
