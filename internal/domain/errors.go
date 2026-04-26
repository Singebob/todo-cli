package domain

import "errors"

var (
	ErrTitleRequired   = errors.New("title is required")
	ErrContextRequired = errors.New("context is required")
	ErrNotFound        = errors.New("todo not found")
	ErrContextNotFound = errors.New("context not found")
	ErrContextExists   = errors.New("context already exists")
	ErrContextHasTodos = errors.New("context still has todos")
)
