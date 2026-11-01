package githubproject

import (
	"context"
	"fmt"

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

	for _, githubProject := range fetchedGithubProjects {
		err = uc.repoPersist.UpStore(ctx, &githubProject)
		if err != nil {
			return fmt.Errorf("GithubProjectUseCase - RefreshMirror - uc.repoPersist.UpStore: %w", err)
		}
	}

	return nil
}
