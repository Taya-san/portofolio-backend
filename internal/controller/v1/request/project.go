package request

import "time"

type (
	CreateProject struct {
		ProjectUrl  string    `json:"project_url" validate:"required"`
		Title       string    `json:"title" validate:"required"`
		Description string    `json:"description" validate:"required"`
		Languages   []string  `json:"languages" validate:"required"`
		CreatedAt   time.Time `json:"createdAt" validate:"required"`
		Visible     bool      `json:"visible" validate:"required"`
		Body        string    `json:"body" validate:"required"`
	} // @name v1.CreateProject

	UpdateProject struct {
		ProjectUrl  string    `json:"project_url" validate:"omitempty"`
		Title       string    `json:"title" validate:"omitempty"`
		Description string    `json:"description" validate:"omitempty"`
		Languages   []string  `json:"languages" validate:"omitempty"`
		CreatedAt   time.Time `json:"createdAt" validate:"omitempty"`
		Visible     bool      `json:"visible" validate:"omitempty"`
		Body        string    `json:"body" validate:"omitempty"`
	} // @name v1.UpdateProject
)
