package project

import (
	"context"
	"time"

	"github.com/Taya-san/portofolio-backend/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../mocks_usecase_project_test.go -package=usecase_test

type (
	Project interface {
		Create(ctx context.Context, projectItem ProjectItem) (entity.Project, error)
		Get(ctx context.Context, projectID string) (entity.Project, error)
		List(ctx context.Context, all bool) ([]entity.Project, int, error)
		Update(ctx context.Context, projectID string, projectItem ProjectItem) (entity.Project, error)
		Delete(ctx context.Context, projectID string) error
		SyncAll(ctx context.Context) error
		SyncByProjectID(ctx context.Context, projectID string) error
	}

	ProjectItem struct {
		ProjectUrl  string
		Title       string
		Description string
		Languages   []string
		CreatedAt   time.Time
		Visible     bool
		Custom      bool
		Body        string
	}
)
