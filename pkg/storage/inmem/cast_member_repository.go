package inmem

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/selmison/code-micro-videos/pkg/cast_member"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

var (
	castMembersOnce sync.Once
	repo            *castMemberRepository
)

// castMemberRepository keeps castMembers in the memory.
// Use it in tests or for development/demo purposes.
type castMemberRepository struct {
	castMembers     map[string]cast_member.CastMember
	castMembersOnce sync.Once
	mu              sync.RWMutex
}

// NewCastMemberRepository returns a new in-memory castMember castMemberRepository.
func NewCastMemberRepository() cast_member.Repository {
	castMembersOnce.Do(func() {
		repo = &castMemberRepository{}
		repo.init()
	})
	return repo
}

func (r *castMemberRepository) init() {
	r.castMembersOnce.Do(func() {
		r.castMembers = make(map[string]cast_member.CastMember)
	})
}

// Store stores an castMember.
func (r *castMemberRepository) Store(_ context.Context, castMember cast_member.CastMember) error {
	r.init()
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, value := range r.castMembers {
		if value.Id() == castMember.Id() {
			return fmt.Errorf("id '%s' %w", castMember.Id(), logger.ErrAlreadyExists)
		}
	}
	r.castMembers[castMember.Id()] = castMember
	return nil
}

// GetAll returns all castMembers.
func (r *castMemberRepository) GetAll(_ context.Context) ([]cast_member.CastMember, error) {
	r.init()
	r.mu.RLock()
	defer r.mu.RUnlock()
	castMembers := make([]cast_member.CastMember, len(r.castMembers))
	// This makes sure castMembers are always returned in the same, sorted order
	keys := make([]string, 0, len(r.castMembers))
	for k := range r.castMembers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for i, key := range keys {
		castMembers[i] = r.castMembers[key]
	}
	return castMembers, nil
}

// GetMany returns a list of castMembers filtered by ids.
func (r *castMemberRepository) GetMany(_ context.Context, ids []string) ([]cast_member.CastMember, error) {
	r.init()
	r.mu.RLock()
	defer r.mu.RUnlock()
	length := len(ids)
	var castMembers []cast_member.CastMember
	var keys []string
	if length > 0 {
		castMembers = make([]cast_member.CastMember, length)
		keys = make([]string, 0, length)
		for _, k := range ids {
			keys = append(keys, k)
		}
	} else {
		length = len(r.castMembers)
		castMembers = make([]cast_member.CastMember, length)
		// This makes sure castMembers are always returned in the same, sorted order
		keys = make([]string, 0, length)
		for k := range r.castMembers {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, key := range keys {
		if compare(r.castMembers[key], nil) {
			castMembers = append(castMembers, r.castMembers[key])
		}
	}
	return castMembers, nil
}

// DeleteAll deletes all castMembers from the castMemberRepository.
func (r *castMemberRepository) DeleteAll(_ context.Context) error {
	r.init()
	r.mu.Lock()
	defer r.mu.Unlock()
	r.castMembers = make(map[string]cast_member.CastMember)
	return nil
}

// GetOne returns a single castMember by its Id.
func (r *castMemberRepository) GetOne(_ context.Context, id string) (cast_member.CastMember, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	castMember, ok := r.castMembers[id]
	if !ok {
		return castMember, fmt.Errorf("%s: %w", id, logger.ErrNotFound)
	}
	return castMember, nil
}

// DeleteOne deletes a single castMember by its Id.
func (r *castMemberRepository) DeleteOne(_ context.Context, id string) error {
	r.init()
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.castMembers[id]
	if !ok {
		return fmt.Errorf("%s: %w", id, logger.ErrNotFound)
	}
	delete(r.castMembers, id)
	return nil
}

// UpdateOne updates a single castMember by its Id.
func (r *castMemberRepository) UpdateOne(_ context.Context, id string, updateCastMember cast_member.UpdateCastMemberDTO) error {
	r.init()
	r.mu.Lock()
	defer r.mu.Unlock()
	castMember, ok := r.castMembers[id]
	if !ok {
		return fmt.Errorf("%s: %w", id, logger.ErrNotFound)
	}
	castMember.Update(updateCastMember)
	r.castMembers[id] = castMember
	return nil
}

func compare(a, b cast_member.CastMember) bool {
	if &a == &b {
		return true
	}
	if a.Id() != b.Id() {
		return false
	}
	if a.Name() != b.Name() {
		return false
	}
	if a.Type() != b.Type() {
		return false
	}
	return true
}
