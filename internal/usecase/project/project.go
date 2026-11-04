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

func New(r project.ProjectRepo, g githubproject.GithubProjectRepo) Project {
	return newTraced(&UseCase{
		repo:       r,
		githubRepo: g,
	})
}

func (uc *UseCase) Create(ctx context.Context, projectItem *ProjectItem) (entity.Project, error) {
	prj := entity.Project{
		ID:          uuid.New().String(),
		ProjectUrl:  projectItem.ProjectUrl,
		Title:       projectItem.Title,
		Languages:   projectItem.Languages,
		Description: projectItem.Description,
		CreatedAt:   projectItem.CreatedAt,
		Visible:     projectItem.Visible,
		Custom:      true,
		Body:        projectItem.Body,
	}

	err := uc.repo.Store(ctx, &prj)
	if err != nil {
		return entity.Project{}, fmt.Errorf("ProjectUseCase - Create - uc.repo.Store: %w", err)
	}

	return prj, nil
}

func (uc *UseCase) Get(ctx context.Context, projectID string) (entity.Project, error) {
	prj, err := uc.repo.GetByID(ctx, projectID)
	if err != nil {
		return entity.Project{}, fmt.Errorf("ProjectUseCase - Get - uc.repo.GetByID: %w", err)
	}

	return prj, nil
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

func (uc *UseCase) Update(ctx context.Context, projectID string, projectItem *ProjectItem) (entity.Project, error) {
	if strings.TrimSpace(projectItem.Title) == "" {
		return entity.Project{}, entity.ErrTitleInvalid
	}

	prj, err := uc.repo.GetByID(
		ctx,
		projectID,
	)
	if err != nil {
		return entity.Project{}, fmt.Errorf("ProjectUseCase - Update - uc.repo.GetById: %w", err)
	}

	prj.ProjectUrl = projectItem.ProjectUrl
	prj.Title = projectItem.Title
	prj.Description = projectItem.Description
	prj.Languages = projectItem.Languages
	prj.Visible = projectItem.Visible
	prj.Custom = true
	prj.Body = projectItem.Body

	err = uc.repo.Update(ctx, &prj)
	if err != nil {
		return entity.Project{}, fmt.Errorf("ProjectUseCase - Update - uc.repo.Update: %w", err)
	}

	return prj, nil
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

	for idx, _ := range githubProjects {
		err := uc.repo.UpStore(ctx, &entity.Project{
			ID:          uuid.NewString(),
			ProjectUrl:  githubProjects[idx].ProjectUrl,
			Title:       githubProjects[idx].Title,
			Description: githubProjects[idx].Description,
			Languages:   githubProjects[idx].Languages,
			CreatedAt:   githubProjects[idx].CreatedAt,
			Visible:     true,
			Custom:      false,
			Body:        githubProjects[idx].Body,
		})
		if err != nil {
			return fmt.Errorf("ProjectUseCase - SyncAll - uc.repo.UpStore: %w", err)
		}
	}

	return nil
}

func (uc *UseCase) SyncByProjectID(ctx context.Context, projectID string) error {
	prj, err := uc.repo.GetByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("ProjectUseCase - SyncByProjectID - uc.repo.GetByID: %w", err)
	}

	githubProject, err := uc.githubRepo.GetByProjectUrl(ctx, prj.ProjectUrl)
	if err != nil {
		return fmt.Errorf("ProjectUseCase - SyncAll - uc.githubRepo.GetByProjectUrl: %w", err)
	}

	err = uc.repo.Update(ctx, &entity.Project{
		ID:          projectID,
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

	return nil
}
