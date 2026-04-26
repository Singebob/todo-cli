package domain

import "github.com/google/uuid"

type TodoRepository interface {
	FindAll() ([]Todo, error)
	Save(todo Todo) error
	Delete(id uuid.UUID) error
}

type ContextRepository interface {
	FindAll() ([]string, error)
	Save(name string) error
	Delete(name string) error
}
