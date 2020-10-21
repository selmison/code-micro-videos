package video

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

// InMemoryStore keeps videos in the memory.
// Use it in tests or for development/demo purposes.
type InMemoryStore struct {
	videos     map[string]Video
	videosOnce sync.Once
	mu         sync.RWMutex
}

// NewInMemoryStore returns a new in-memory video store.
func NewInMemoryStore() Repository {
	store := &InMemoryStore{}
	store.init()
	return store
}

func (s *InMemoryStore) init() {
	s.videosOnce.Do(func() {
		s.videos = make(map[string]Video)
	})
}

// Store stores an video.
func (s *InMemoryStore) Store(_ context.Context, video Video) error {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, v := range s.videos {
		if v.Id == video.Id {
			return fmt.Errorf("id '%s' %w", video.Id, logger.ErrAlreadyExists)
		}
	}
	s.videos[video.Id] = video
	return nil
}

// GetAll returns all videos.
func (s *InMemoryStore) GetAll(_ context.Context) ([]Video, error) {
	s.init()
	s.mu.RLock()
	defer s.mu.RUnlock()
	videos := make([]Video, len(s.videos))
	// This makes sure videos are always returned in the same, sorted order
	keys := make([]string, 0, len(s.videos))
	for k := range s.videos {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, key := range keys {
		videos[i] = s.videos[key]
	}
	return videos, nil
}

// GetMany returns a list of videos filtered by ids.
func (s *InMemoryStore) GetMany(_ context.Context, ids []string) ([]Video, error) {
	s.init()
	s.mu.RLock()
	defer s.mu.RUnlock()
	length := len(ids)
	var videos []Video
	var keys []string
	if length > 0 {
		videos = make([]Video, length)
		keys = make([]string, 0, length)
		for _, k := range ids {
			keys = append(keys, k)
		}
	} else {
		length = len(s.videos)
		videos = make([]Video, length)
		// This makes sure videos are always returned in the same, sorted order
		keys = make([]string, 0, length)
		for k := range s.videos {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, key := range keys {
		if compare(s.videos[key], Video{}) {
			videos = append(videos, s.videos[key])
		}
	}
	return videos, nil
}

// DeleteAll deletes all videos from the store.
func (s *InMemoryStore) DeleteAll(_ context.Context) error {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.videos = make(map[string]Video)
	return nil
}

// GetOne returns a single video by its Id.
func (s *InMemoryStore) GetOne(_ context.Context, id string) (Video, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	video, ok := s.videos[id]
	if !ok {
		return video, fmt.Errorf("%s: %w", id, logger.ErrNotFound)
	}
	return video, nil
}

// DeleteOne deletes a single video by its Id.
func (s *InMemoryStore) DeleteOne(_ context.Context, id string) error {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.videos[id]
	if !ok {
		return fmt.Errorf("%s: %w", id, logger.ErrNotFound)
	}
	delete(s.videos, id)
	return nil
}

// UpdateOne updates a single video by its Id.
func (s *InMemoryStore) UpdateOne(_ context.Context, id string, updateVideo UpdateVideo) error {
	s.init()
	s.mu.Lock()
	defer s.mu.Unlock()
	video, ok := s.videos[id]
	if !ok {
		return fmt.Errorf("%s: %w", id, logger.ErrNotFound)
	}
	updateVideo.update(video)
	s.videos[id] = video
	return nil
}

func compare(a Video, b Video) bool {
	if &a == &b {
		return true
	}
	if a.Id != b.Id {
		return false
	}
	if a.Title != b.Title {
		return false
	}
	if a.Description != b.Description {
		return false
	}
	if a.YearLaunched != b.YearLaunched {
		return false
	}
	if a.Opened != b.Opened {
		return false
	}

	if a.Rating != b.Rating {
		return false
	}

	if a.Duration != b.Duration {
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

	if len(a.CategoriesId) != len(b.CategoriesId) {
		return false
	}
	for i, v := range a.CategoriesId {
		if b.CategoriesId[i] != v {
			return false
		}
	}
	if a.VideoFileHandler == nil && b.VideoFileHandler != nil || a.VideoFileHandler != nil && b.VideoFileHandler == nil {
		return false
	}
	//TODO: add deep compare to VideoFileHandler
	return true
}
