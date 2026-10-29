package entity

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrProjectNotFound    = errors.New("task not found")
	ErrProjectForbidden   = errors.New("unauthorized task activity")
	ErrTitleInvalid       = errors.New("project title cannot be blank")
)
