package project

import (
	"context"
	"fmt"
	"strings"

	"github.com/Taya-san/portofolio-backend/internal/entity"
	githubproject "github.com/Taya-san/portofolio-backend/internal/repo/persist/github_project"
	"github.com/Taya-san/portofolio-backend/internal/repo/persist/project"
	"github.com/google/uuid"
)

type UseCase struct {
	repo       project.ProjectRepo
	githubRepo githubproject.GithubProjectRepo
}

func New(r project.ProjectRepo) Project {
	return newTraced(&UseCase{repo: r})
}

func (uc *UseCase) Create(ctx context.Context, projectItem ProjectItem) (entity.Project, error) {

	project := entity.Project{
		ID:         uuid.New().String(),
		ProjectUrl: projectItem.ProjectUrl,
		Title:      projectItem.Title,
		Languages:  projectItem.Languages,
		CreatedAt:  projectItem.CreatedAt,
		Visible:    projectItem.Visible,
		Custom:     true,
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
	project.Description = projectItem.Description
	project.Languages = projectItem.Languages
	project.Visible = projectItem.Visible
	project.Custom = true
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

func (uc *UseCase) SyncAll(ctx context.Context) error {
	githubProjects, _, err := uc.githubRepo.List(ctx)
	if err != nil {
		return fmt.Errorf("ProjectUseCase - SyncAll - uc.githubRepo.List: %w", err)
	}

	for _, githubProject := range githubProjects {
		err := uc.repo.UpStore(ctx, &entity.Project{
			ID:          uuid.NewString(),
			ProjectUrl:  githubProject.ProjectUrl,
			Title:       githubProject.Title,
			Description: githubProject.Description,
			Languages:   githubProject.Languages,
			CreatedAt:   githubProject.CreatedAt,
			Visible:     true,
			Custom:      false,
			Body:        githubProject.Body,
		})
		if err != nil {
			return fmt.Errorf("ProjectUseCase - SyncAll - uc.repo.UpStore: %w", err)
		}
	}

	return nil
}

func (uc *UseCase) SyncByProjectID(ctx context.Context, projectID string) error {
	project, err := uc.repo.GetByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("ProjectUseCase - SyncByProjectID - uc.repo.GetByID: %w", err)
	}

	githubProject, err := uc.githubRepo.GetByProjectUrl(ctx, project.ProjectUrl)
	if err != nil {
		return fmt.Errorf("ProjectUseCase - SyncAll - uc.githubRepo.GetByProjectUrl: %w", err)
	}

	err = uc.repo.Update(ctx, &entity.Project{
		ID:          uuid.NewString(),
		ProjectUrl:  githubProject.ProjectUrl,
		Title:       githubProject.Title,
		Description: githubProject.Description,
		Languages:   githubProject.Languages,
		CreatedAt:   githubProject.CreatedAt,
		Visible:     project.Visible,
		Custom:      false,
		Body:        githubProject.Body,
	})
	if err != nil {
		return fmt.Errorf("ProjectUseCase - SyncAll - uc.repo.UpStore: %w", err)
	}

	return nil
}
