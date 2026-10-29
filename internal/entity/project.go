package entity

import "time"

type Project struct {
	ID         string    `json:"id"`
	ProjectUrl string    `json:"project_url"`
	Title      string    `json:"title"`
	Languages  []string  `json:"languages"`
	CreatedAt  time.Time `json:"createdAt"`
	Visible    bool      `json:"visible"`
	Custom     bool      `json:"custom"`
	Body       string    `json:"body"`
} // @name entity.Project
