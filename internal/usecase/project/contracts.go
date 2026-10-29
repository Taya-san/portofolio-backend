package project

import (
	"context"

	"github.com/Taya-san/portofolio-backend/internal/entity"
)

type (
	Project interface {
		Create(ctx context.Context, projectItem ProjectItem) (entity.Project, error)
		Get(ctx context.Context, projectID string) (entity.Project, error)
		List(ctx context.Context, all bool) ([]entity.Project, int, error)
		Update(ctx context.Context, projectID string, projectItem ProjectItem) (entity.Project, error)
		Delete(ctx context.Context, projectID string) error
	}

	ProjectItem struct {
		ProjectUrl string
		Title      string
		Languages  []string
		Visible    bool
		Custom     bool
		Body       string
	}
)
