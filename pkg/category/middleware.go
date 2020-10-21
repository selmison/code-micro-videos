package category

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

func (m validationMiddleware) Create(ctx context.Context, newCategory NewCategory) (Category, error) {
	if newCategory.isEmpty() {
		return Category{}, errors.WithStack(fmt.Errorf("'%s' param %w", "newCategory", logger.ErrCouldNotBeEmpty))
	}
	newCategory.Name = strings.ToLower(strings.TrimSpace(newCategory.Name))
	if newCategory.Name == "" {
		return Category{}, errors.WithStack(fmt.Errorf("'%s' field %w", "Name", logger.ErrIsRequired))
	}
	return m.next.Create(ctx, newCategory)
}

func (m validationMiddleware) Destroy(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		//
		return errors.WithStack(fmt.Errorf("'%s' params %w", "id", logger.ErrCouldNotBeEmpty))
	}
	return m.next.Destroy(ctx, id)
}

func (m validationMiddleware) List(ctx context.Context) ([]Category, error) {
	return m.next.List(ctx)
}

func (m validationMiddleware) Show(ctx context.Context, id string) (Category, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Category{}, errors.WithStack(fmt.Errorf("'%s' param %w", "id", logger.ErrCouldNotBeEmpty))
	}
	return m.next.Show(ctx, id)
}

func (m validationMiddleware) Update(ctx context.Context, id string, updateCategory UpdateCategory) error {
	if updateCategory.isEmpty() {
		return nil
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.WithStack(fmt.Errorf("'%s' params %w", "id", logger.ErrCouldNotBeEmpty))
	}
	if updateCategory.Name != nil && strings.TrimSpace(*updateCategory.Name) == "" {
		return fmt.Errorf("the %s field %w", "Name", logger.ErrCouldNotBeEmpty)
	}
	return m.next.Update(ctx, id, updateCategory)
}

//func (m validationMiddleware) List(ctx context.Context, opts ListUsersOptions) ([]*User, *ListMeta, error) {
//	if err := opts.Validate(); err != nil {
//		return nil, nil, err
//	}
//	return m.next.ListUsers(ctx, opts)
//}
