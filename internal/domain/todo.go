package domain

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID          uuid.UUID
	Title       string
	Comment     string
	Context     string
	Category    string
	Done        bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

func NewTodo(title, comment, context, category string) (Todo, error) {
	if title == "" {
		return Todo{}, ErrTitleRequired
	}
	if context == "" {
		return Todo{}, ErrContextRequired
	}
	return Todo{
		ID:        uuid.New(),
		Title:     title,
		Comment:   comment,
		Context:   context,
		Category:  category,
		Done:      false,
		CreatedAt: time.Now(),
	}, nil
}

func (t *Todo) ToggleDone() {
	t.Done = !t.Done
	if t.Done {
		now := time.Now()
		t.CompletedAt = &now
	} else {
		t.CompletedAt = nil
	}
}
