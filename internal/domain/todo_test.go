package domain_test

import (
	"testing"

	"todo-cli/internal/domain"
)

func TestNewTodo(t *testing.T) {
	todo, err := domain.NewTodo("Buy milk", "at the store", "perso", "courses")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if todo.Title != "Buy milk" {
		t.Errorf("title = %q, want %q", todo.Title, "Buy milk")
	}
	if todo.Comment != "at the store" {
		t.Errorf("comment = %q, want %q", todo.Comment, "at the store")
	}
	if todo.Context != "perso" {
		t.Errorf("context = %q, want %q", todo.Context, "perso")
	}
	if todo.Category != "courses" {
		t.Errorf("category = %q, want %q", todo.Category, "courses")
	}
	if todo.Done {
		t.Error("new todo should not be done")
	}
	if todo.CompletedAt != nil {
		t.Error("new todo should have nil CompletedAt")
	}
	if todo.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
	if todo.ID == [16]byte{} {
		t.Error("ID should be generated")
	}
}

func TestNewTodo_EmptyTitle(t *testing.T) {
	_, err := domain.NewTodo("", "comment", "perso", "")
	if err != domain.ErrTitleRequired {
		t.Errorf("err = %v, want ErrTitleRequired", err)
	}
}

func TestNewTodo_EmptyContext(t *testing.T) {
	_, err := domain.NewTodo("title", "comment", "", "")
	if err != domain.ErrContextRequired {
		t.Errorf("err = %v, want ErrContextRequired", err)
	}
}

func TestNewTodo_OptionalFields(t *testing.T) {
	todo, err := domain.NewTodo("title", "", "perso", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if todo.Comment != "" {
		t.Errorf("comment should be empty, got %q", todo.Comment)
	}
	if todo.Category != "" {
		t.Errorf("category should be empty, got %q", todo.Category)
	}
}

func TestToggleDone(t *testing.T) {
	todo, _ := domain.NewTodo("task", "", "perso", "")

	todo.ToggleDone()
	if !todo.Done {
		t.Error("todo should be done after first toggle")
	}
	if todo.CompletedAt == nil {
		t.Error("CompletedAt should be set when done")
	}

	todo.ToggleDone()
	if todo.Done {
		t.Error("todo should not be done after second toggle")
	}
	if todo.CompletedAt != nil {
		t.Error("CompletedAt should be nil when undone")
	}
}
