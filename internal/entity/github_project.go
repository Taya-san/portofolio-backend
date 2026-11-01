package entity

import "time"

type GithubProject struct {
	ID          string
	ProjectUrl  string
	Title       string
	Description string
	Languages   []string
	CreatedAt   time.Time
	Body        string
} // @name entity.GithubProject
