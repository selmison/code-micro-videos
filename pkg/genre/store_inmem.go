package genre

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

// InMemoryStore keeps genres in the memory.
// Use it in tests or for development/demo purposes.
type InMemoryStore struct {
	genres     map[string]Genre
	genresOnce sync.Once
	mu         sync.RWMutex
}

// NewInMemoryStore returns a new in-memory genre store.
func NewInMemoryStore() Repository {
	store := &InMemoryStore{}
	store.init()
	return store
}

func (s *InMemoryStore) init() {
	s.genresOnce.Do(func() {
		s.genres = make(map[string]Genre)
	})
}

// Store stores an genre.
func (s *InMemoryStore) Store(_ context.Context, genre Genre) error {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, value := range s.genres {
		if value.Name == genre.Name {
			return fmt.Errorf("name '%s' %w", genre.Name, logger.ErrAlreadyExists)
		}
	}
	s.genres[genre.Id] = genre
	return nil
}

// GetAll returns all genres.
func (s *InMemoryStore) GetAll(_ context.Context) ([]Genre, error) {
	s.init()
	s.mu.RLock()
	defer s.mu.RUnlock()
	genres := make([]Genre, len(s.genres))
	// This makes sure genres are always returned in the same, sorted order
	keys := make([]string, 0, len(s.genres))
	for k := range s.genres {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, key := range keys {
		genres[i] = s.genres[key]
	}
	return genres, nil
}

// GetMany returns a list of genres filtered by ids.
func (s *InMemoryStore) GetMany(_ context.Context, ids []string) ([]Genre, error) {
	s.init()
	s.mu.RLock()
	defer s.mu.RUnlock()
	length := len(ids)
	var genres []Genre
	var keys []string
	if length > 0 {
		genres = make([]Genre, length)
		keys = make([]string, 0, length)
		for _, k := range ids {
			keys = append(keys, k)
		}
	} else {
		length = len(s.genres)
		genres = make([]Genre, length)
		// This makes sure genres are always returned in the same, sorted order
		keys = make([]string, 0, length)
		for k := range s.genres {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, key := range keys {
		if compare(s.genres[key], Genre{}) {
			genres = append(genres, s.genres[key])
		}
	}
	return genres, nil
}

// DeleteAll deletes all genres from the store.
func (s *InMemoryStore) DeleteAll(_ context.Context) error {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.genres = make(map[string]Genre)
	return nil
}

// GetOne returns a single genre by its Id.
func (s *InMemoryStore) GetOne(_ context.Context, id string) (Genre, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	genre, ok := s.genres[id]
	if !ok {
		return genre, fmt.Errorf("%s: %w", id, logger.ErrNotFound)
	}
	return genre, nil
}

// DeleteOne deletes a single genre by its Id.
func (s *InMemoryStore) DeleteOne(_ context.Context, id string) error {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.genres[id]
	if !ok {
		return fmt.Errorf("%s: %w", id, logger.ErrNotFound)
	}
	delete(s.genres, id)
	return nil
}

// UpdateOne updates a single genre by its Id.
func (s *InMemoryStore) UpdateOne(_ context.Context, id string, updateGenre UpdateGenre) error {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()
	genre, ok := s.genres[id]
	if !ok {
		return fmt.Errorf("%s: %w", id, logger.ErrNotFound)
	}
	updateGenre.update(genre)
	s.genres[id] = genre
	return nil
}

func compare(a, b Genre) bool {
	if &a == &b {
		return true
	}
	if a.Id != b.Id {
		return false
	}
	if a.Name != b.Name {
		return false
	}
	if len(a.CategoriesId) != len(b.CategoriesId) {
		return false
	}
	for i, v := range a.CategoriesId {
		if b.CategoriesId[i] != v {
			return false
		}
	}
	return true
}
