package inmemory

import (
	"sort"

	"github.com/google/uuid"
	"todo-cli/internal/domain"
)

type TodoRepository struct {
	todos map[uuid.UUID]domain.Todo
}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{todos: make(map[uuid.UUID]domain.Todo)}
}

func (r *TodoRepository) FindAll() ([]domain.Todo, error) {
	todos := make([]domain.Todo, 0, len(r.todos))
	for _, t := range r.todos {
		todos = append(todos, t)
	}
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].CreatedAt.Before(todos[j].CreatedAt)
	})
	return todos, nil
}

func (r *TodoRepository) Save(todo domain.Todo) error {
	r.todos[todo.ID] = todo
	return nil
}

func (r *TodoRepository) Delete(id uuid.UUID) error {
	if _, ok := r.todos[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.todos, id)
	return nil
}
