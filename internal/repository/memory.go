package repository

import (
	"errors"
	"sync"

	"moremail/email-finder/internal/domain"
)

type MemoryRepository struct {
	data map[string]*domain.EmailAnalysis
	mu   sync.Mutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		data: make(map[string]*domain.EmailAnalysis),
	}
}

func (r *MemoryRepository) Save(a *domain.EmailAnalysis) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[a.ID] = a
}

func (r *MemoryRepository) FindByID(id string) (*domain.EmailAnalysis, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	a, ok := r.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return a, nil
}
