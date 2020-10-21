package cast_member

import (
	"context"
	"fmt"
	"strings"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

// ValidationMiddleware describes a service (as opposed to endpoint) middleware.
type ValidationMiddleware func(Service) Service

type validationMiddleware struct {
	next Service
}

// NewValidationMiddleware returns a service ValidationMiddleware.
func NewValidationMiddleware() ValidationMiddleware {
	return func(next Service) Service {
		return validationMiddleware{next: next}
	}
}

func (m validationMiddleware) Create(ctx context.Context, newCastMember NewCastMemberDTO) (CastMember, error) {
	return m.next.Create(ctx, newCastMember)
}

func (m validationMiddleware) Destroy(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("'%s' param %w", "id", logger.ErrCouldNotBeEmpty)
	}
	return m.next.Destroy(ctx, id)
}

func (m validationMiddleware) List(ctx context.Context) ([]CastMember, error) {
	return m.next.List(ctx)
}

func (m validationMiddleware) Show(ctx context.Context, id string) (CastMember, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("'%s' param %w", "id", logger.ErrCouldNotBeEmpty)
	}
	return m.next.Show(ctx, id)
}

func (m validationMiddleware) Update(ctx context.Context, id string, updateCastMember UpdateCastMemberDTO) error {
	if updateCastMember.isEmpty() {
		return nil
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("'%s' param %w", "id", logger.ErrCouldNotBeEmpty)
	}
	if err := updateCastMember.validate(); err != nil {
		return err
	}
	return m.next.Update(ctx, id, updateCastMember)
}
