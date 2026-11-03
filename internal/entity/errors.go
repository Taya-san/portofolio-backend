package entity

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrProjectNotFound    = errors.New("project not found")
	ErrProjectForbidden   = errors.New("unauthorized project activity")
	ErrTitleInvalid       = errors.New("project title cannot be blank")
)
