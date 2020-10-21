package category

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

// InMemoryStore keeps categories in the memory.
// Use it in tests or for development/demo purposes.
type InMemoryStore struct {
	categories     map[string]Category
	categoriesOnce sync.Once
	mu             sync.RWMutex
}

// NewInMemoryStore returns a new in-memory category store.
func NewInMemoryStore() Repository {
	store := &InMemoryStore{}
	store.init()
	return store
}

func (s *InMemoryStore) init() {
	s.categoriesOnce.Do(func() {
		s.categories = make(map[string]Category)
	})
}

// Store stores an category.
func (s *InMemoryStore) Store(_ context.Context, category Category) error {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, value := range s.categories {
		if value.Name == category.Name {
			return fmt.Errorf("name '%s' %w", category.Name, logger.ErrAlreadyExists)
		}
	}
	s.categories[category.Id] = category
	return nil
}

// GetAll returns all categories.
func (s *InMemoryStore) GetAll(_ context.Context) ([]Category, error) {
	s.init()
	s.mu.RLock()
	defer s.mu.RUnlock()
	categories := make([]Category, len(s.categories))
	// This makes sure categories are always returned in the same, sorted order
	keys := make([]string, 0, len(s.categories))
	for k := range s.categories {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, key := range keys {
		categories[i] = s.categories[key]
	}
	return categories, nil
}

// GetMany returns a list of categories filtered by ids.
func (s *InMemoryStore) GetMany(_ context.Context, ids []string) ([]Category, error) {
	s.init()
	s.mu.RLock()
	defer s.mu.RUnlock()
	length := len(ids)
	var categories []Category
	var keys []string
	if length > 0 {
		categories = make([]Category, length)
		keys = make([]string, 0, length)
		for _, k := range ids {
			keys = append(keys, k)
		}
	} else {
		length = len(s.categories)
		categories = make([]Category, length)
		// This makes sure categories are always returned in the same, sorted order
		keys = make([]string, 0, length)
		for k := range s.categories {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, key := range keys {
		if compare(s.categories[key], Category{}) {
			categories = append(categories, s.categories[key])
		}
	}
	return categories, nil
}

// DeleteAll deletes all categories from the store.
func (s *InMemoryStore) DeleteAll(_ context.Context) error {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.categories = make(map[string]Category)
	return nil
}

// GetOne returns a single category by its Id.
func (s *InMemoryStore) GetOne(_ context.Context, id string) (Category, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	category, ok := s.categories[id]
	if !ok {
		return category, fmt.Errorf("%s: %w", id, logger.ErrNotFound)
	}
	return category, nil
}

// DeleteOne deletes a single category by its Id.
func (s *InMemoryStore) DeleteOne(_ context.Context, id string) error {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.categories[id]
	if !ok {
		return fmt.Errorf("%s: %w", id, logger.ErrNotFound)
	}
	delete(s.categories, id)
	return nil
}

// UpdateOne updates a single category by its Id.
func (s *InMemoryStore) UpdateOne(_ context.Context, id string, updateCategory UpdateCategory) error {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()
	category, ok := s.categories[id]
	if !ok {
		return fmt.Errorf("%s: %w", id, logger.ErrNotFound)
	}
	updateCategory.update(category)
	s.categories[id] = category
	return nil
}

func compare(a, b Category) bool {
	if &a == &b {
		return true
	}
	if a.Id != b.Id {
		return false
	}
	if a.Name != b.Name {
		return false
	}
	if a.Description != b.Description {
		return false
	}
	if len(a.GenresId) != len(b.GenresId) {
		return false
	}
	for i, v := range a.GenresId {
		if b.GenresId[i] != v {
			return false
		}
	}
	return true
}
