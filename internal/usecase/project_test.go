package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Taya-san/portofolio-backend/internal/entity"
	persistProjectRepo "github.com/Taya-san/portofolio-backend/internal/repo/persist/project"
	"github.com/Taya-san/portofolio-backend/internal/usecase/project"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newProjectUseCase(t *testing.T) (project.Project, *MockProjectRepo, *MockGithubProjectRepo) {
	t.Helper()

	ctrl := gomock.NewController(t)

	mockRepo := NewMockProjectRepo(ctrl)
	mockGithubRepo := NewMockGithubProjectRepo(ctrl)
	useCase := project.New(mockRepo, mockGithubRepo)

	return useCase, mockRepo, mockGithubRepo
}

func TestProjectCreate(t *testing.T) {
	t.Parallel()

	t.Run("create success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newProjectUseCase(t)
		mockRepo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(nil)

		timeNow := time.Now().UTC()

		p2, err := uc.Create(context.Background(),
			&project.ProjectItem{
				ProjectUrl:  "project-url-something.com",
				Title:       "the TITLE OF THE PROJECT",
				Description: "the DESCRIPTION OF THE PROJECT",
				Languages:   []string{"English", "Rust", "Jawascript"},
				CreatedAt:   timeNow,
				Visible:     true,
				Body:        "the BODY OF THE PROJECT",
			},
		)

		require.NoError(t, err)
		assert.NotEmpty(t, p2.ID)
		assert.Equal(t, "the TITLE OF THE PROJECT", p2.Title)
		assert.Equal(t, timeNow, p2.CreatedAt)
	})
}

func TestProjectGet(t *testing.T) {
	t.Parallel()

	timeNow := time.Now().UTC()

	expectedProject := entity.Project{
		ID:          "project-id-098",
		ProjectUrl:  "project-url-something.com",
		Title:       "the TITLE OF THE PROJECT",
		Description: "the DESCRIPTION OF THE PROJECT",
		Languages:   []string{"English", "Rust", "Jawascript"},
		CreatedAt:   timeNow,
		Visible:     true,
		Custom:      true,
		Body:        "the BODY OF THE PROJECT",
	}

	t.Run("get success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newProjectUseCase(t)
		mockRepo.EXPECT().GetByID(gomock.Any(), "project-id-098").Return(expectedProject, nil)

		p2, err := uc.Get(context.Background(), "project-id-098")

		require.NoError(t, err)
		assert.Equal(t, expectedProject, p2)
	})
}

func TestProjectList(t *testing.T) {
	t.Parallel()

	timeNow := time.Now().UTC()

	project1 := entity.Project{
		ID:          "project-id-098",
		ProjectUrl:  "project-url-something.com",
		Title:       "the TITLE OF THE PROJECT",
		Description: "the DESCRIPTION OF THE PROJECT",
		Languages:   []string{"English", "Rust", "Jawascript"},
		CreatedAt:   timeNow,
		Visible:     true,
		Custom:      true,
		Body:        "the BODY OF THE PROJECT",
	}

	project2 := entity.Project{
		ID:          "project-id-098-2",
		ProjectUrl:  "project-url-something.com",
		Title:       "the TITLE OF THE PROJECT 2",
		Description: "the DESCRIPTION OF THE PROJECT 2",
		Languages:   []string{"English", "Rust", "Jawascript"},
		CreatedAt:   timeNow,
		Visible:     true,
		Custom:      true,
		Body:        "the BODY OF THE PROJECT 2 ",
	}

	t.Run("get success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newProjectUseCase(t)
		mockRepo.EXPECT().List(gomock.Any(), persistProjectRepo.ProjectFilter{All: true}).Return([]entity.Project{project1, project2}, 2, nil)

		projects, n, err := uc.List(context.Background(), true)

		require.NoError(t, err)
		assert.Equal(t, 2, n)
		assert.Len(t, projects, 2)
	})
}

func TestProjectUpdate(t *testing.T) {
	t.Parallel()

	t.Run("create success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newProjectUseCase(t)

		timeNow := time.Now().UTC()

		existingProject := entity.Project{
			ID:          "project-id-098",
			ProjectUrl:  "project-url-something.com",
			Title:       "the TITLE OF THE PROJECT",
			Description: "the DESCRIPTION OF THE PROJECT",
			Languages:   []string{"English", "Rust", "Jawascript"},
			CreatedAt:   timeNow,
			Visible:     true,
			Custom:      true,
			Body:        "the BODY OF THE PROJECT",
		}

		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), "project-id-098").Return(existingProject, nil)

		p2, err := uc.Update(context.Background(), "project-id-098",
			&project.ProjectItem{
				ProjectUrl:  "project-url-something.com",
				Title:       "the TITLE OF THE PROJECT",
				Description: "the DESCRIPTION OF THE PROJECT",
				Languages:   []string{"English", "Rust", "Jawascript"},
				CreatedAt:   timeNow,
				Visible:     true,
				Body:        "the BODY OF THE PROJECT",
			},
		)

		require.NoError(t, err)
		assert.NotEmpty(t, p2.ID)
		assert.Equal(t, "the TITLE OF THE PROJECT", p2.Title)
		assert.Equal(t, timeNow, p2.CreatedAt)
	})
}

func TestProjectDelete(t *testing.T) {
	t.Parallel()

	t.Run("create success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newProjectUseCase(t)
		mockRepo.EXPECT().Delete(gomock.Any(), "project-id-098").Return(nil)

		err := uc.Delete(context.Background(), "project-id-098")

		require.NoError(t, err)
	})

	t.Run("delete not found", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, _ := newProjectUseCase(t)
		mockRepo.EXPECT().Delete(gomock.Any(), "project-id-missing").Return(entity.ErrProjectNotFound)

		err := uc.Delete(context.Background(), "project-id-missing")

		require.ErrorIs(t, err, entity.ErrProjectNotFound)
	})
}

func TestProjectSyncAll(t *testing.T) {
	t.Parallel()

	timeNow := time.Now().UTC()

	var githubProjects []entity.GithubProject

	githubProject := entity.GithubProject{
		ID:          "project-id-098",
		ProjectUrl:  "project-url-something.com",
		Title:       "the TITLE OF THE PROJECT",
		Description: "the DESCRIPTION OF THE PROJECT",
		Languages:   []string{"English", "Rust", "Jawascript"},
		CreatedAt:   timeNow,
		Body:        "the BODY OF THE PROJECT",
	}

	githubProjects = append(githubProjects, githubProject)

	t.Run("sync all success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, mockGithubRepo := newProjectUseCase(t)
		mockGithubRepo.EXPECT().List(gomock.Any()).Return(githubProjects, 1, nil)
		mockRepo.EXPECT().UpStore(gomock.Any(), gomock.Any()).Return(nil)

		err := uc.SyncAll(context.Background())

		assert.NoError(t, err)
	})
}

func TestProjectSyncByProjectID(t *testing.T) {
	t.Parallel()

	timeNow := time.Now().UTC()

	var githubProjects []entity.GithubProject

	githubProject := entity.GithubProject{
		ID:          "project-id-098",
		ProjectUrl:  "project-url-something-new.com",
		Title:       "the TITLE OF THE PROJECT",
		Description: "the DESCRIPTION OF THE PROJECT",
		Languages:   []string{"English", "Rust", "Jawascript"},
		CreatedAt:   timeNow,
		Body:        "the BODY OF THE PROJECT",
	}

	githubProjects = append(githubProjects, githubProject)

	prj := entity.Project{
		ID:          "project-id-098",
		ProjectUrl:  "project-url-something.com",
		Title:       "the TITLE OF THE PROJECT",
		Description: "the DESCRIPTION OF THE PROJECT",
		Languages:   []string{"English", "Rust", "Jawascript"},
		CreatedAt:   timeNow,
		Visible:     true,
		Custom:      true,
		Body:        "the BODY OF THE PROJECT",
	}

	updateProject := &entity.Project{
		ID:          "project-id-098",
		ProjectUrl:  githubProject.ProjectUrl,
		Title:       githubProject.Title,
		Description: githubProject.Description,
		Languages:   githubProject.Languages,
		CreatedAt:   githubProject.CreatedAt,
		Visible:     true,
		Custom:      false,
		Body:        githubProject.Body,
	}

	t.Run("sync by id success", func(t *testing.T) {
		t.Parallel()

		uc, mockRepo, mockGithubRepo := newProjectUseCase(t)
		mockGithubRepo.EXPECT().GetByProjectUrl(gomock.Any(), prj.ProjectUrl).Return(githubProject, nil)
		mockRepo.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(prj, nil)
		mockRepo.EXPECT().Update(gomock.Any(), updateProject).Return(nil)

		err := uc.SyncByProjectID(context.Background(), "project-id-098")

		assert.NoError(t, err)
	})
}

func TestProjectCreate_RepoError(t *testing.T) {
	t.Parallel()

	uc, mockRepo, _ := newProjectUseCase(t)
	mockRepo.EXPECT().Store(gomock.Any(), gomock.Any()).Return(errRepoGeneric)

	timeNow := time.Now().UTC()

	_, err := uc.Create(context.Background(),
		&project.ProjectItem{
			ProjectUrl:  "project-url-something.com",
			Title:       "the TITLE OF THE PROJECT",
			Description: "the DESCRIPTION OF THE PROJECT",
			Languages:   []string{"English", "Rust", "Jawascript"},
			CreatedAt:   timeNow,
			Visible:     true,
			Body:        "the BODY OF THE PROJECT",
		},
	)

	require.Error(t, err)
	require.ErrorIs(t, err, errRepoGeneric)
}
