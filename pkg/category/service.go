package category

import (
	"context"

	"github.com/go-kit/kit/log"

	"github.com/selmison/code-micro-videos/pkg/id_generator"
)

// Service manages categories.
type Service interface {
	// Create creates a new category.
	Create(ctx context.Context, newCategory NewCategory) (Category, error)

	// Destroy destroys a category.
	Destroy(ctx context.Context, id string) error

	// List returns a list of categories.
	List(ctx context.Context) ([]Category, error)

	// Show returns the details of a category.
	Show(ctx context.Context, id string) (Category, error)

	// Update updates an existing category.
	Update(ctx context.Context, id string, updateCategory UpdateCategory) error
}

// Category represents a single Category.
type Category struct {
	Id          string
	Name        string
	Description string
	GenresId    []string
	IsValidated bool
}

// NewCategory contains the details of a new Category.
type NewCategory struct {
	Name        string
	Description string
	GenresId    []string
	IsValidated bool
}

func (c NewCategory) toCategory(id string) Category {
	return Category{
		Id:          id,
		Name:        c.Name,
		Description: c.Description,
		GenresId:    c.GenresId,
		IsValidated: c.IsValidated,
	}
}

// UpdateCategory contains updates of an existing category.
type UpdateCategory struct {
	Name        *string
	Description *string
	GenresId    []string
	IsValidated *bool
}

func (c UpdateCategory) update(category Category) Category {
	if c.Name != nil {
		category.Name = *c.Name
	}
	if c.Description != nil {
		category.Description = *c.Description
	}
	if c.GenresId != nil {
		category.GenresId = c.GenresId
	}
	if c.IsValidated != nil {
		category.IsValidated = *c.IsValidated
	}
	return category
}

type service struct {
	idGenerator id_generator.IdGenerator
	repo        Repository
}

// NewService returns a new Service with all of the expected middlewares wired in.
func NewService(idGenerator id_generator.IdGenerator, repo Repository, logger log.Logger) Service {
	var svc Service
	{
		svc = service{idGenerator: idGenerator, repo: repo}
		svc = NewValidationMiddleware()(svc)
		svc = NewLoggingMiddleware(logger)(svc)
	}
	return svc
}

func (svc service) Create(ctx context.Context, newCategory NewCategory) (Category, error) {
	id, err := svc.idGenerator.Generate()
	if err != nil {
		return Category{}, err
	}
	category := newCategory.toCategory(id)
	err = svc.repo.Store(ctx, category)
	if err != nil {
		return Category{}, err
	}
	return category, nil
}

func (svc service) Destroy(ctx context.Context, id string) error {
	return svc.repo.DeleteOne(ctx, id)
}

func (svc service) List(ctx context.Context) ([]Category, error) {
	categories, err := svc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (svc service) Show(ctx context.Context, id string) (Category, error) {
	category, err := svc.repo.GetOne(ctx, id)
	if err != nil {
		return Category{}, err
	}
	return category, nil
}

func (svc service) Update(ctx context.Context, id string, updateCategory UpdateCategory) error {
	err := svc.repo.UpdateOne(ctx, id, updateCategory)
	if err != nil {
		return err
	}
	return nil
}

func (c NewCategory) isEmpty() bool {
	return c.compare(NewCategory{})
}

func (c NewCategory) compare(b NewCategory) bool {
	if &c == &b {
		return true
	}
	if c.Name != b.Name {
		return false
	}
	if c.Description != b.Description {
		return false
	}
	if len(c.GenresId) != len(b.GenresId) {
		return false
	}
	for i, v := range c.GenresId {
		if b.GenresId[i] != v {
			return false
		}
	}
	return true
}

func (c UpdateCategory) isEmpty() bool {
	return c.compare(UpdateCategory{})
}

func (c UpdateCategory) compare(b UpdateCategory) bool {
	if &c == &b {
		return true
	}
	if c.Name != b.Name {
		return false
	}
	if c.Description != b.Description {
		return false
	}
	if len(c.GenresId) != len(b.GenresId) {
		return false
	}
	for i, v := range c.GenresId {
		if b.GenresId[i] != v {
			return false
		}
	}
	return true
}
