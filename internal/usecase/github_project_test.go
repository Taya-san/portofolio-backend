package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Taya-san/portofolio-backend/internal/entity"
	githubproject "github.com/Taya-san/portofolio-backend/internal/usecase/github_project"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var errRepoGeneric = errors.New("repository error")

func newGithubProjectUseCase(t *testing.T) (githubproject.GithubProject, *MockGithubProjectRepo, *MockGithubProjectApi) {
	t.Helper()

	ctrl := gomock.NewController(t)

	mockRepo := NewMockGithubProjectRepo(ctrl)
	mockApi := NewMockGithubProjectApi(ctrl)
	useCase := githubproject.New(mockRepo, mockApi)

	return useCase, mockRepo, mockApi
}

func TestProjectRefreshMirror(t *testing.T) {
	t.Parallel()

	timeNow := time.Now().UTC()

	t.Run("refresh success", func(t *testing.T) {
		t.Parallel()

		var fetchedProjects []entity.GithubProject

		fetchedProject := entity.GithubProject{
			ID:          "project-id-098",
			ProjectUrl:  "project-url-something.com",
			Title:       "the TITLE OF THE PROJECT",
			Description: "the DESCRIPTION OF THE PROJECT",
			Languages:   []string{"English", "Rust", "Jawascript"},
			CreatedAt:   timeNow,
			Body:        "the BODY OF THE PROJECT",
		}

		fetchedProjects = append(fetchedProjects, fetchedProject)

		uc, mockRepo, mockApi := newGithubProjectUseCase(t)
		mockRepo.EXPECT().UpStore(gomock.Any(), gomock.Any()).Return(nil)
		mockApi.EXPECT().Fetch(gomock.Any()).Return(fetchedProjects, nil)

		err := uc.RefreshMirror(context.Background())

		assert.NoError(t, err)
	})
}

func TestProjectGetByProjectUrl(t *testing.T) {
	t.Parallel()

	timeNow := time.Now().UTC()

	t.Run("get success", func(t *testing.T) {
		t.Parallel()

		fetchedProject := entity.GithubProject{
			ID:          "project-id-098",
			ProjectUrl:  "project-url-something.com",
			Title:       "the TITLE OF THE PROJECT",
			Description: "the DESCRIPTION OF THE PROJECT",
			Languages:   []string{"English", "Rust", "Jawascript"},
			CreatedAt:   timeNow,
			Body:        "the BODY OF THE PROJECT",
		}

		uc, mockRepo, _ := newGithubProjectUseCase(t)
		mockRepo.EXPECT().GetByProjectUrl(gomock.Any(), gomock.Any()).Return(fetchedProject, nil)

		p2, err := uc.GetByProjectUrl(context.Background(), "project-url-something.com")

		require.NoError(t, err)
		assert.Equal(t, fetchedProject.ProjectUrl, p2.ProjectUrl)
		assert.Equal(t, fetchedProject.Body, p2.Body)
	})
}
