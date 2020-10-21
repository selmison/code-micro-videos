package video

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

func (m validationMiddleware) Create(ctx context.Context, newVideo NewVideo) (Video, error) {
	if newVideo.isEmpty() {
		return Video{}, errors.WithStack(fmt.Errorf("'%s' param %w", "newVideo", logger.ErrCouldNotBeEmpty))
	}
	if err := newVideo.Validate(); err != nil {
		return Video{}, errors.WithStack(err)
	}
	return m.next.Create(ctx, newVideo)
}

func (m validationMiddleware) Destroy(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.WithStack(fmt.Errorf("'%s' %w", "id", logger.ErrCouldNotBeEmpty))
	}
	return m.next.Destroy(ctx, id)
}

func (m validationMiddleware) List(ctx context.Context) ([]Video, error) {
	return m.next.List(ctx)
}

func (m validationMiddleware) Show(ctx context.Context, id string) (Video, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Video{}, errors.WithStack(fmt.Errorf("'%s' field %w", "id", logger.ErrCouldNotBeEmpty))
	}
	return m.next.Show(ctx, id)
}

func (m validationMiddleware) Update(ctx context.Context, id string, updateVideo UpdateVideo) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.WithStack(fmt.Errorf("'%s' field %w", "id", logger.ErrCouldNotBeEmpty))
	}
	if err := updateVideo.Validate(); err != nil {
		return errors.WithStack(err)
	}
	return m.next.Update(ctx, id, updateVideo)
}

//func (m validationMiddleware) List(ctx context.Context, opts ListUsersOptions) ([]*User, *ListMeta, error) {
//	if err := opts.Validate(); err != nil {
//		return nil, nil, err
//	}
//	return m.next.ListUsers(ctx, opts)
//}
