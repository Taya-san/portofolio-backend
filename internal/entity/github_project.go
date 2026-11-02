package entity

import "time"

type GithubProject struct {
	ID          string    `json:"id"`
	ProjectUrl  string    `json:"project_url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Languages   []string  `json:"languages"`
	CreatedAt   time.Time `json:"createdAt"`
	Body        string    `json:"body"`
} // @name entity.GithubProject
