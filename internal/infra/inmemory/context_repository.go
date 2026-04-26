package inmemory

import "todo-cli/internal/domain"

type ContextRepository struct {
	contexts []string
}

func NewContextRepository() *ContextRepository {
	return &ContextRepository{}
}

func (r *ContextRepository) FindAll() ([]string, error) {
	result := make([]string, len(r.contexts))
	copy(result, r.contexts)
	return result, nil
}

func (r *ContextRepository) Save(name string) error {
	for _, c := range r.contexts {
		if c == name {
			return domain.ErrContextExists
		}
	}
	r.contexts = append(r.contexts, name)
	return nil
}

func (r *ContextRepository) Delete(name string) error {
	for i, c := range r.contexts {
		if c == name {
			r.contexts = append(r.contexts[:i], r.contexts[i+1:]...)
			return nil
		}
	}
	return domain.ErrContextNotFound
}
