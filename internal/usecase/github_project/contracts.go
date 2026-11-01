package githubproject

import (
	"context"
)

//go:generate mockgen -source=contracts.go -destination=../mocks_usecase_githubproject_test.go -package=usecase_test

type (
	GithubProject interface {
		RefreshMirror(ctx context.Context) error
	}
)
