package project

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Taya-san/portofolio-backend/internal/entity"
	"github.com/Taya-san/portofolio-backend/internal/repo/project"
	"github.com/google/uuid"
)

type UseCase struct {
	repo project.ProjectRepo
}

func New(r project.ProjectRepo) Project {
	return newTraced(&UseCase{repo: r})
}

func (uc *UseCase) Create(ctx context.Context, projectItem ProjectItem) (entity.Project, error) {
	now := time.Now().UTC()

	project := entity.Project{
		ID:         uuid.New().String(),
		ProjectUrl: projectItem.ProjectUrl,
		Title:      projectItem.Title,
		Languages:  projectItem.Languages,
		CreatedAt:  now,
		Visible:    projectItem.Visible,
		Custom:     projectItem.Custom,
		Body:       projectItem.Body,
	}

	err := uc.repo.Store(ctx, &project)
	if err != nil {
		return entity.Project{}, fmt.Errorf("ProjectUseCase - Create - uc.repo.Store: %w", err)
	}

	return project, nil
}

func (uc *UseCase) Get(ctx context.Context, projectID string) (entity.Project, error) {
	project, err := uc.repo.GetByID(ctx, projectID)
	if err != nil {
		return entity.Project{}, fmt.Errorf("ProjectUseCase - Get - uc.repo.GetByID: %w", err)
	}

	return project, nil
}

func (uc *UseCase) List(ctx context.Context, all bool) ([]entity.Project, int, error) {
	projects, total, err := uc.repo.List(
		ctx,
		project.ProjectFilter{All: all},
	)
	if err != nil {
		return nil, 0, fmt.Errorf("ProjectUseCase - List - uc.repo.List: %w", err)
	}

	return projects, total, nil
}

func (uc *UseCase) Update(ctx context.Context, projectID string, projectItem ProjectItem) (entity.Project, error) {
	if strings.TrimSpace(projectItem.Title) == "" {
		return entity.Project{}, entity.ErrTitleInvalid
	}

	project, err := uc.repo.GetByID(
		ctx,
		projectID,
	)
	if err != nil {
		return entity.Project{}, fmt.Errorf("ProjectUseCase - Update - uc.repo.GetById: %w", err)
	}

	project.ProjectUrl = projectItem.ProjectUrl
	project.Title = projectItem.Title
	project.Languages = projectItem.Languages
	project.Visible = projectItem.Visible
	project.Body = projectItem.Body

	err = uc.repo.Update(ctx, &project)
	if err != nil {
		return entity.Project{}, fmt.Errorf("ProjectUseCase - Update - uc.repo.Update: %w", err)
	}

	return project, nil
}

func (uc *UseCase) Delete(ctx context.Context, projectID string) error {
	err := uc.repo.Delete(ctx, projectID)
	if err != nil {
		return fmt.Errorf("ProjectUseCase - Delete - uc.repo.Delete: %w", err)
	}

	return nil
}
