package githubproject

import (
	"context"
	"fmt"

	"github.com/Taya-san/portofolio-backend/internal/entity"
	githubprojectpersist "github.com/Taya-san/portofolio-backend/internal/repo/persist/github_project"
	githubprojectwebapi "github.com/Taya-san/portofolio-backend/internal/repo/webapi/github_project"
)

type UseCase struct {
	repoPersist githubprojectpersist.GithubProjectRepo
	webApi      githubprojectwebapi.GithubProjectApi
}

func New(r githubprojectpersist.GithubProjectRepo, w githubprojectwebapi.GithubProjectApi) GithubProject {
	return newTraced(&UseCase{
		repoPersist: r,
		webApi:      w,
	})
}

func (uc *UseCase) RefreshMirror(ctx context.Context) error {
	fetchedGithubProjects, err := uc.webApi.Fetch(ctx)
	if err != nil {
		return fmt.Errorf("GithubProjectUseCase - RefreshMirror - uc.webApi.Fetch: %w", err)
	}

	for idx := range fetchedGithubProjects {
		err = uc.repoPersist.UpStore(ctx, &fetchedGithubProjects[idx])
		if err != nil {
			return fmt.Errorf("GithubProjectUseCase - RefreshMirror - uc.repoPersist.UpStore: %w", err)
		}
	}

	return nil
}

func (uc *UseCase) GetByProjectUrl(ctx context.Context, projectUrl string) (entity.GithubProject, error) {
	githubProject, err := uc.repoPersist.GetByProjectUrl(ctx, projectUrl)
	if err != nil {
		return entity.GithubProject{}, fmt.Errorf("GithubProjectUseCase - GetByProjectUrl - uc.repoPersist.GetByProjectUrl: %w", err)
	}

	return githubProject, err
}
