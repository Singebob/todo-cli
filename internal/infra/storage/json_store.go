package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"todo-cli/internal/domain"
)

type todoDTO struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Comment     string     `json:"comment"`
	Context     string     `json:"context"`
	Category    string     `json:"category"`
	Done        bool       `json:"done"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

func toDTO(t domain.Todo) todoDTO {
	return todoDTO{
		ID:          t.ID,
		Title:       t.Title,
		Comment:     t.Comment,
		Context:     t.Context,
		Category:    t.Category,
		Done:        t.Done,
		CreatedAt:   t.CreatedAt,
		CompletedAt: t.CompletedAt,
	}
}

func toDomain(d todoDTO) domain.Todo {
	return domain.Todo{
		ID:          d.ID,
		Title:       d.Title,
		Comment:     d.Comment,
		Context:     d.Context,
		Category:    d.Category,
		Done:        d.Done,
		CreatedAt:   d.CreatedAt,
		CompletedAt: d.CompletedAt,
	}
}

type JSONStore struct {
	dir      string
	mu       sync.RWMutex
	todos    map[uuid.UUID]domain.Todo
	contexts []string
}

func NewJSONStore() (*JSONStore, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(home, ".todo-cli")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	s := &JSONStore{
		dir:   dir,
		todos: make(map[uuid.UUID]domain.Todo),
	}
	if err := s.loadTodos(); err != nil {
		return nil, err
	}
	if err := s.loadContexts(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *JSONStore) todosPath() string    { return filepath.Join(s.dir, "todos.json") }
func (s *JSONStore) contextsPath() string { return filepath.Join(s.dir, "contexts.json") }

func (s *JSONStore) loadTodos() error {
	data, err := os.ReadFile(s.todosPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}
	var dtos []todoDTO
	if err := json.Unmarshal(data, &dtos); err != nil {
		return err
	}
	for _, d := range dtos {
		t := toDomain(d)
		s.todos[t.ID] = t
	}
	return nil
}

func (s *JSONStore) loadContexts() error {
	data, err := os.ReadFile(s.contextsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &s.contexts)
}

func (s *JSONStore) flushTodos() error {
	todos := s.sortedTodos()
	dtos := make([]todoDTO, len(todos))
	for i, t := range todos {
		dtos[i] = toDTO(t)
	}
	return s.writeJSON(s.todosPath(), dtos)
}

func (s *JSONStore) flushContexts() error {
	return s.writeJSON(s.contextsPath(), s.contexts)
}

func (s *JSONStore) writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (s *JSONStore) sortedTodos() []domain.Todo {
	todos := make([]domain.Todo, 0, len(s.todos))
	for _, t := range s.todos {
		todos = append(todos, t)
	}
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].CreatedAt.Before(todos[j].CreatedAt)
	})
	return todos
}

// --- TodoRepository ---

func (s *JSONStore) FindAll() ([]domain.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sortedTodos(), nil
}

func (s *JSONStore) Save(todo domain.Todo) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.todos[todo.ID] = todo
	return s.flushTodos()
}

func (s *JSONStore) Delete(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		return domain.ErrNotFound
	}
	delete(s.todos, id)
	return s.flushTodos()
}

// --- ContextRepository (via ContextStore adapter) ---

type ContextStore struct {
	store *JSONStore
}

func (s *JSONStore) ContextStore() *ContextStore {
	return &ContextStore{store: s}
}

func (cs *ContextStore) FindAll() ([]string, error) {
	cs.store.mu.RLock()
	defer cs.store.mu.RUnlock()
	result := make([]string, len(cs.store.contexts))
	copy(result, cs.store.contexts)
	return result, nil
}

func (cs *ContextStore) Save(name string) error {
	cs.store.mu.Lock()
	defer cs.store.mu.Unlock()
	for _, c := range cs.store.contexts {
		if c == name {
			return domain.ErrContextExists
		}
	}
	cs.store.contexts = append(cs.store.contexts, name)
	return cs.store.flushContexts()
}

func (cs *ContextStore) Delete(name string) error {
	cs.store.mu.Lock()
	defer cs.store.mu.Unlock()
	for i, c := range cs.store.contexts {
		if c == name {
			cs.store.contexts = append(cs.store.contexts[:i], cs.store.contexts[i+1:]...)
			return cs.store.flushContexts()
		}
	}
	return domain.ErrContextNotFound
}
