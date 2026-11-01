package githubproject

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"
	"uuid"

	"github.com/Taya-san/portofolio-backend/internal/entity"
)

type userProject struct {
	Title       string    `json:"name"`
	FullName    string    `json:"full_name"`
	ProjectUrl  string    `json:"html_url"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type lang struct {
	name  string
	bytes int
}

type Repo struct {
	*http.Client
	GithubToken string
}

func New(h *http.Client, GithubToken string) GithubProjectApi {
	return newTraced(&Repo{
		h,
		GithubToken,
	})
}

func (r *Repo) Fetch(ctx context.Context) ([]entity.GithubProject, error) {
	allUserProjects, err := r.fetchUserProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("GithubProjectApi - Fetch - r.fetchUserProjects: %w", err)
	}

	allUserProjectsData, err := r.fetchProjectData(ctx, allUserProjects)
	if err != nil {
		return nil, fmt.Errorf("GithubProjectApi - Fetch - r.fetchProjectData: %w", err)
	}

	return allUserProjectsData, nil
}

func (r *Repo) fetchUserProjects(ctx context.Context) ([]userProject, error) {
	var allUserProjects []userProject

	for page := 1; ; page++ {
		var userProjects []userProject

		url := fmt.Sprintf(
			"https://api.github.com/user/repos?per_page=100&page=%d",
			page,
		)

		req, err := http.NewRequestWithContext(ctx,
			http.MethodGet,
			url,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("GithubProjectApi - fetchUserProjects - http.NewRequestWithContext: %w", err)
		}

		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Authorization", "Bearer "+r.GithubToken)
		req.Header.Set("X-GitHub-Api-Version", "2026-03-10")

		resp, err := r.Do(req)
		if err != nil {
			return nil, fmt.Errorf("GithubProjectApi - fetchUserProjects - r.Do: %w", err)
		}
		header := resp.Header.Get("link")

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("GithubProjectApi - fetchUserProjects - r.Do: %s", resp.Status)
		}

		if err := json.UnmarshalRead(resp.Body, &userProjects); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("GithubProjectApi - fetchUserProjects - json.UnmarshalRead: %w", err)
		}

		resp.Body.Close()

		allUserProjects = append(allUserProjects, userProjects...)

		if !strings.Contains(header, `rel="next"`) {
			break
		}
	}

	return allUserProjects, nil
}

func (r *Repo) fetchProjectData(ctx context.Context, allUserProjects []userProject) ([]entity.GithubProject, error) {
	allProjects := make([]entity.GithubProject, 0, len(allUserProjects))
	var fetchedLangs map[string]int

	for _, userProject := range allUserProjects {
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/readme",
			userProject.FullName,
		)

		req, err := http.NewRequestWithContext(ctx,
			http.MethodGet,
			url,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("GithubProjectApi - fetchProjectData - http.NewRequestWithContext: %w", err)
		}

		req.Header.Set("Accept", "application/vnd.github.raw+json")
		req.Header.Set("Authorization", "Bearer "+r.GithubToken)
		req.Header.Set("X-GitHub-Api-Version", "2026-03-10")

		resp, err := r.Do(req)
		if err != nil {
			return nil, fmt.Errorf("GithubProjectApi - fetchProjectData - r.Do: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("GithubProjectApi - fetchProjectData - r.Do: %s", resp.Status)
		}

		projectBodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("GithubProjectApi - fetchProjectData - io.ReadAll: %w", err)
		}

		projectBody := string(projectBodyBytes)

		url = fmt.Sprintf(
			"https://api.github.com/repos/%s/languages",
			userProject.FullName,
		)

		req, err = http.NewRequestWithContext(ctx,
			http.MethodGet,
			url,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("GithubProjectApi - fetchProjectData - http.NewRequestWithContext: %w", err)
		}

		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Authorization", "Bearer "+r.GithubToken)
		req.Header.Set("X-GitHub-Api-Version", "2026-03-10")

		resp, err = r.Do(req)
		if err != nil {
			return nil, fmt.Errorf("GithubProjectApi - fetchProjectData - r.Do: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("GithubProjectApi - fetchProjectData - r.Do: %s", resp.Status)
		}

		if err := json.UnmarshalRead(resp.Body, &fetchedLangs); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("GithubProjectApi - fetchProjectData - json.UnmarshalRead: %w", err)
		}
		resp.Body.Close()

		langs := make([]lang, 0, len(fetchedLangs))

		for k, v := range fetchedLangs {
			langs = append(langs,
				lang{
					name:  k,
					bytes: v,
				},
			)
		}

		slices.SortFunc(langs,
			func(a lang, b lang) int {
				return b.bytes - a.bytes
			},
		)

		langNames := make([]string, 0, len(langs))
		for _, lang := range langs {
			langNames = append(langNames, lang.name)
		}

		allProjects = append(allProjects, entity.GithubProject{
			ID:          uuid.New().String(),
			Title:       userProject.Title,
			ProjectUrl:  userProject.ProjectUrl,
			Description: userProject.Description,
			Languages:   langNames,
			CreatedAt:   userProject.CreatedAt,
			Body:        projectBody,
		})
	}

	return allProjects, nil
}
