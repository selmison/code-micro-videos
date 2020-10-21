package genre

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

type validationMiddleware struct {
	next Service
}

// NewValidationMiddleware returns a service Middleware.
func NewValidationMiddleware() Middleware {
	return func(next Service) Service {
		return validationMiddleware{next: next}
	}
}

func (m validationMiddleware) Create(ctx context.Context, newGenre NewGenre) (Genre, error) {
	if newGenre.isEmpty() {
		return Genre{}, errors.WithStack(fmt.Errorf("'%s' param %w", "newGenre", logger.ErrCouldNotBeEmpty))
	}
	newGenre.Name = strings.ToLower(strings.TrimSpace(newGenre.Name))
	if newGenre.Name == "" {
		return Genre{}, errors.WithStack(fmt.Errorf("'%s' field %w", "Name", logger.ErrIsRequired))
	}
	return m.next.Create(ctx, newGenre)
}

func (m validationMiddleware) Destroy(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.WithStack(fmt.Errorf("'%s' param %w", "id", logger.ErrCouldNotBeEmpty))
	}
	return m.next.Destroy(ctx, id)
}

func (m validationMiddleware) List(ctx context.Context) ([]Genre, error) {
	return m.next.List(ctx)
}

func (m validationMiddleware) Show(ctx context.Context, id string) (Genre, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Genre{}, errors.WithStack(fmt.Errorf("'%s' param %w", "id", logger.ErrCouldNotBeEmpty))
	}
	return m.next.Show(ctx, id)
}

func (m validationMiddleware) Update(ctx context.Context, id string, updateGenre UpdateGenre) error {
	if updateGenre.isEmpty() {
		return nil
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.WithStack(fmt.Errorf("'%s' param %w", "id", logger.ErrCouldNotBeEmpty))
	}
	if updateGenre.Name != nil && strings.TrimSpace(*updateGenre.Name) == "" {
		return fmt.Errorf("the %s field %w", "Name", logger.ErrCouldNotBeEmpty)
	}
	return m.next.Update(ctx, id, updateGenre)
}

//func (m validationMiddleware) List(ctx context.Context, opts ListUsersOptions) ([]*User, *ListMeta, error) {
//	if err := opts.Validate(); err != nil {
//		return nil, nil, err
//	}
//	return m.next.ListUsers(ctx, opts)
//}
