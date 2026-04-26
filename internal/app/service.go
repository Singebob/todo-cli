package app

import (
	"sort"

	"github.com/google/uuid"
	"todo-cli/internal/domain"
)

type TodoService struct {
	todos    domain.TodoRepository
	contexts domain.ContextRepository
}

func NewTodoService(todos domain.TodoRepository, contexts domain.ContextRepository) *TodoService {
	return &TodoService{todos: todos, contexts: contexts}
}

func (s *TodoService) ListByContext(ctx string) ([]domain.Todo, error) {
	all, err := s.todos.FindAll()
	if err != nil {
		return nil, err
	}
	var result []domain.Todo
	for _, t := range all {
		if t.Context == ctx && !t.Done {
			result = append(result, t)
		}
	}
	return result, nil
}

func (s *TodoService) ListCompletedByContext(ctx string) ([]domain.Todo, error) {
	all, err := s.todos.FindAll()
	if err != nil {
		return nil, err
	}
	var result []domain.Todo
	for _, t := range all {
		if t.Context == ctx && t.Done {
			result = append(result, t)
		}
	}
	return result, nil
}

func (s *TodoService) ListByContextAndCategory(ctx, category string) ([]domain.Todo, error) {
	all, err := s.todos.FindAll()
	if err != nil {
		return nil, err
	}
	var result []domain.Todo
	for _, t := range all {
		if t.Context == ctx && t.Category == category && !t.Done {
			result = append(result, t)
		}
	}
	return result, nil
}

func (s *TodoService) GetTodo(id uuid.UUID) (domain.Todo, error) {
	all, err := s.todos.FindAll()
	if err != nil {
		return domain.Todo{}, err
	}
	for _, t := range all {
		if t.ID == id {
			return t, nil
		}
	}
	return domain.Todo{}, domain.ErrNotFound
}

func (s *TodoService) CreateTodo(title, comment, context, category string) (domain.Todo, error) {
	if err := s.contextMustExist(context); err != nil {
		return domain.Todo{}, err
	}
	todo, err := domain.NewTodo(title, comment, context, category)
	if err != nil {
		return domain.Todo{}, err
	}
	if err := s.todos.Save(todo); err != nil {
		return domain.Todo{}, err
	}
	return todo, nil
}

func (s *TodoService) UpdateTodo(id uuid.UUID, title, comment, context, category string) (domain.Todo, error) {
	if err := s.contextMustExist(context); err != nil {
		return domain.Todo{}, err
	}
	todo, err := s.GetTodo(id)
	if err != nil {
		return domain.Todo{}, err
	}
	if title == "" {
		return domain.Todo{}, domain.ErrTitleRequired
	}
	todo.Title = title
	todo.Comment = comment
	todo.Context = context
	todo.Category = category
	if err := s.todos.Save(todo); err != nil {
		return domain.Todo{}, err
	}
	return todo, nil
}

func (s *TodoService) ToggleDone(id uuid.UUID) (domain.Todo, error) {
	todo, err := s.GetTodo(id)
	if err != nil {
		return domain.Todo{}, err
	}
	todo.ToggleDone()
	if err := s.todos.Save(todo); err != nil {
		return domain.Todo{}, err
	}
	return todo, nil
}

func (s *TodoService) DeleteTodo(id uuid.UUID) error {
	return s.todos.Delete(id)
}

func (s *TodoService) AllContexts() ([]string, error) {
	return s.contexts.FindAll()
}

func (s *TodoService) CreateContext(name string) error {
	if name == "" {
		return domain.ErrContextRequired
	}
	return s.contexts.Save(name)
}

func (s *TodoService) DeleteContext(name string) error {
	all, err := s.todos.FindAll()
	if err != nil {
		return err
	}
	for _, t := range all {
		if t.Context == name {
			return domain.ErrContextHasTodos
		}
	}
	return s.contexts.Delete(name)
}

func (s *TodoService) CategoriesForContext(ctx string) ([]string, error) {
	all, err := s.todos.FindAll()
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{})
	for _, t := range all {
		if t.Context == ctx && t.Category != "" {
			seen[t.Category] = struct{}{}
		}
	}
	categories := make([]string, 0, len(seen))
	for c := range seen {
		categories = append(categories, c)
	}
	sort.Strings(categories)
	return categories, nil
}

func (s *TodoService) contextMustExist(name string) error {
	contexts, err := s.contexts.FindAll()
	if err != nil {
		return err
	}
	for _, c := range contexts {
		if c == name {
			return nil
		}
	}
	return domain.ErrContextNotFound
}
