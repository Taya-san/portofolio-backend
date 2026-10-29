package request

import "time"

type (
	CreateProject struct {
		ProjectUrl string    `json:"project_url"`
		Title      string    `json:"title"`
		Languages  []string  `json:"languages"`
		CreatedAt  time.Time `json:"createdAt"`
		Visible    bool      `json:"visible"`
		Custom     bool      `json:"custom"`
		Body       string    `json:"body"`
	} // @name v1.CreateProject

	UpdateProject struct {
		ProjectUrl string    `json:"project_url" validate:"omitempty"`
		Title      string    `json:"title" validate:"omitempty"`
		Languages  []string  `json:"languages" validate:"omitempty"`
		CreatedAt  time.Time `json:"createdAt" validate:"omitempty"`
		Visible    bool      `json:"visible" validate:"omitempty"`
		Custom     bool      `json:"custom" validate:"omitempty"`
		Body       string    `json:"body" validate:"omitempty"`
	} // @name v1.UpdateProject
)
