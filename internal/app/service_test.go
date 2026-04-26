package app_test

import (
	"testing"

	"github.com/google/uuid"
	"todo-cli/internal/app"
	"todo-cli/internal/domain"
	"todo-cli/internal/infra/inmemory"
)

func setup(t *testing.T) (*app.TodoService, *inmemory.TodoRepository, *inmemory.ContextRepository) {
	t.Helper()
	todos := inmemory.NewTodoRepository()
	contexts := inmemory.NewContextRepository()
	svc := app.NewTodoService(todos, contexts)
	return svc, todos, contexts
}

func setupWithContext(t *testing.T, ctx string) *app.TodoService {
	t.Helper()
	svc, _, contexts := setup(t)
	if err := contexts.Save(ctx); err != nil {
		t.Fatalf("failed to add context: %v", err)
	}
	return svc
}

func TestCreateContext(t *testing.T) {
	svc, _, _ := setup(t)

	if err := svc.CreateContext("perso"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	contexts, _ := svc.AllContexts()
	if len(contexts) != 1 || contexts[0] != "perso" {
		t.Errorf("contexts = %v, want [perso]", contexts)
	}
}

func TestCreateContext_Empty(t *testing.T) {
	svc, _, _ := setup(t)
	if err := svc.CreateContext(""); err != domain.ErrContextRequired {
		t.Errorf("err = %v, want ErrContextRequired", err)
	}
}

func TestCreateContext_Duplicate(t *testing.T) {
	svc, _, _ := setup(t)
	svc.CreateContext("perso")
	if err := svc.CreateContext("perso"); err != domain.ErrContextExists {
		t.Errorf("err = %v, want ErrContextExists", err)
	}
}

func TestDeleteContext(t *testing.T) {
	svc, _, _ := setup(t)
	svc.CreateContext("perso")

	if err := svc.DeleteContext("perso"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	contexts, _ := svc.AllContexts()
	if len(contexts) != 0 {
		t.Errorf("contexts should be empty, got %v", contexts)
	}
}

func TestDeleteContext_NotFound(t *testing.T) {
	svc, _, _ := setup(t)
	if err := svc.DeleteContext("nope"); err != domain.ErrContextNotFound {
		t.Errorf("err = %v, want ErrContextNotFound", err)
	}
}

func TestCreateTodo(t *testing.T) {
	svc := setupWithContext(t, "perso")

	todo, err := svc.CreateTodo("Buy milk", "at the store", "perso", "courses")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if todo.Title != "Buy milk" {
		t.Errorf("title = %q, want %q", todo.Title, "Buy milk")
	}
	if todo.Context != "perso" {
		t.Errorf("context = %q, want %q", todo.Context, "perso")
	}
	if todo.Done {
		t.Error("new todo should not be done")
	}
}

func TestCreateTodo_ContextMustExist(t *testing.T) {
	svc, _, _ := setup(t)

	_, err := svc.CreateTodo("task", "", "unknown", "")
	if err != domain.ErrContextNotFound {
		t.Errorf("err = %v, want ErrContextNotFound", err)
	}
}

func TestCreateTodo_EmptyTitle(t *testing.T) {
	svc := setupWithContext(t, "perso")

	_, err := svc.CreateTodo("", "", "perso", "")
	if err != domain.ErrTitleRequired {
		t.Errorf("err = %v, want ErrTitleRequired", err)
	}
}

func TestListByContext(t *testing.T) {
	svc := setupWithContext(t, "perso")
	svc.CreateContext("work")

	svc.CreateTodo("task1", "", "perso", "")
	svc.CreateTodo("task2", "", "work", "")
	svc.CreateTodo("task3", "", "perso", "")

	todos, err := svc.ListByContext("perso")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(todos) != 2 {
		t.Fatalf("got %d todos, want 2", len(todos))
	}
	if todos[0].Title != "task1" || todos[1].Title != "task3" {
		t.Errorf("todos = [%s, %s], want [task1, task3]", todos[0].Title, todos[1].Title)
	}
}

func TestListByContext_ExcludesCompleted(t *testing.T) {
	svc := setupWithContext(t, "perso")

	todo, _ := svc.CreateTodo("task1", "", "perso", "")
	svc.CreateTodo("task2", "", "perso", "")
	svc.ToggleDone(todo.ID)

	todos, _ := svc.ListByContext("perso")
	if len(todos) != 1 {
		t.Fatalf("got %d todos, want 1", len(todos))
	}
	if todos[0].Title != "task2" {
		t.Errorf("title = %q, want %q", todos[0].Title, "task2")
	}
}

func TestListCompletedByContext(t *testing.T) {
	svc := setupWithContext(t, "perso")

	todo, _ := svc.CreateTodo("task1", "", "perso", "")
	svc.CreateTodo("task2", "", "perso", "")
	svc.ToggleDone(todo.ID)

	completed, _ := svc.ListCompletedByContext("perso")
	if len(completed) != 1 {
		t.Fatalf("got %d completed, want 1", len(completed))
	}
	if completed[0].Title != "task1" {
		t.Errorf("title = %q, want %q", completed[0].Title, "task1")
	}
}

func TestListByContextAndCategory(t *testing.T) {
	svc := setupWithContext(t, "perso")

	svc.CreateTodo("task1", "", "perso", "courses")
	svc.CreateTodo("task2", "", "perso", "admin")
	svc.CreateTodo("task3", "", "perso", "courses")

	todos, _ := svc.ListByContextAndCategory("perso", "courses")
	if len(todos) != 2 {
		t.Fatalf("got %d todos, want 2", len(todos))
	}
}

func TestToggleDone(t *testing.T) {
	svc := setupWithContext(t, "perso")
	todo, _ := svc.CreateTodo("task", "", "perso", "")

	toggled, err := svc.ToggleDone(todo.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !toggled.Done {
		t.Error("todo should be done")
	}
	if toggled.CompletedAt == nil {
		t.Error("CompletedAt should be set")
	}

	untoggled, _ := svc.ToggleDone(todo.ID)
	if untoggled.Done {
		t.Error("todo should not be done after second toggle")
	}
	if untoggled.CompletedAt != nil {
		t.Error("CompletedAt should be nil")
	}
}

func TestUpdateTodo(t *testing.T) {
	svc := setupWithContext(t, "perso")
	todo, _ := svc.CreateTodo("old title", "old comment", "perso", "cat1")

	updated, err := svc.UpdateTodo(todo.ID, "new title", "new comment", "perso", "cat2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Title != "new title" {
		t.Errorf("title = %q, want %q", updated.Title, "new title")
	}
	if updated.Comment != "new comment" {
		t.Errorf("comment = %q, want %q", updated.Comment, "new comment")
	}
	if updated.Category != "cat2" {
		t.Errorf("category = %q, want %q", updated.Category, "cat2")
	}
}

func TestUpdateTodo_ContextMustExist(t *testing.T) {
	svc := setupWithContext(t, "perso")
	todo, _ := svc.CreateTodo("task", "", "perso", "")

	_, err := svc.UpdateTodo(todo.ID, "task", "", "unknown", "")
	if err != domain.ErrContextNotFound {
		t.Errorf("err = %v, want ErrContextNotFound", err)
	}
}

func TestUpdateTodo_EmptyTitle(t *testing.T) {
	svc := setupWithContext(t, "perso")
	todo, _ := svc.CreateTodo("task", "", "perso", "")

	_, err := svc.UpdateTodo(todo.ID, "", "", "perso", "")
	if err != domain.ErrTitleRequired {
		t.Errorf("err = %v, want ErrTitleRequired", err)
	}
}

func TestUpdateTodo_NotFound(t *testing.T) {
	svc := setupWithContext(t, "perso")
	_, err := svc.UpdateTodo(uuid.New(), "title", "", "perso", "")
	if err != domain.ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestDeleteTodo(t *testing.T) {
	svc := setupWithContext(t, "perso")
	todo, _ := svc.CreateTodo("task", "", "perso", "")

	if err := svc.DeleteTodo(todo.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	todos, _ := svc.ListByContext("perso")
	if len(todos) != 0 {
		t.Errorf("todos should be empty, got %d", len(todos))
	}
}

func TestDeleteTodo_NotFound(t *testing.T) {
	svc := setupWithContext(t, "perso")
	if err := svc.DeleteTodo(uuid.New()); err != domain.ErrNotFound {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestGetTodo(t *testing.T) {
	svc := setupWithContext(t, "perso")
	created, _ := svc.CreateTodo("task", "comment", "perso", "cat")

	got, err := svc.GetTodo(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Title != "task" {
		t.Errorf("title = %q, want %q", got.Title, "task")
	}
}

func TestCategoriesForContext(t *testing.T) {
	svc := setupWithContext(t, "perso")
	svc.CreateTodo("t1", "", "perso", "courses")
	svc.CreateTodo("t2", "", "perso", "admin")
	svc.CreateTodo("t3", "", "perso", "courses")
	svc.CreateTodo("t4", "", "perso", "")

	cats, err := svc.CategoriesForContext("perso")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cats) != 2 {
		t.Fatalf("got %d categories, want 2", len(cats))
	}
	if cats[0] != "admin" || cats[1] != "courses" {
		t.Errorf("categories = %v, want [admin, courses]", cats)
	}
}

func TestTodosSortedOldestFirst(t *testing.T) {
	svc := setupWithContext(t, "perso")

	svc.CreateTodo("first", "", "perso", "")
	svc.CreateTodo("second", "", "perso", "")
	svc.CreateTodo("third", "", "perso", "")

	todos, _ := svc.ListByContext("perso")
	if len(todos) != 3 {
		t.Fatalf("got %d todos, want 3", len(todos))
	}
	if todos[0].Title != "first" || todos[1].Title != "second" || todos[2].Title != "third" {
		t.Errorf("order = [%s, %s, %s], want [first, second, third]",
			todos[0].Title, todos[1].Title, todos[2].Title)
	}
}
