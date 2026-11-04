package entity

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrProjectNotFound    = errors.New("project not found")
	ErrTitleInvalid       = errors.New("project title cannot be blank")
	ErrRateLimit          = errors.New("api rate limit exceeded")
	ErrStatusNotOK        = errors.New("http response not OK")
)
