package genre

import (
	"context"

	"github.com/go-kit/kit/log"

	"github.com/selmison/code-micro-videos/pkg/id_generator"
)

type Service interface {
	// Create creates a new genre.
	Create(ctx context.Context, newGenre NewGenre) (Genre, error)

	// Destroy destroys a genre.
	Destroy(ctx context.Context, name string) error

	// List returns a list of genres.
	List(ctx context.Context) ([]Genre, error)

	// Show returns the details of a genre.
	Show(ctx context.Context, id string) (Genre, error)

	// Update updates an existing genre.
	Update(ctx context.Context, id string, updateGenre UpdateGenre) error
}

// Genre represents a single Genre.
type Genre struct {
	Id           string
	Name         string
	CategoriesId []string
	IsValidated  bool
}

// NewGenre contains the details of a new Genre.
type NewGenre struct {
	Name         string
	CategoriesId []string
	IsValidated  bool
}

func (g NewGenre) toGenre(id string) Genre {
	return Genre{
		Id:           id,
		Name:         g.Name,
		CategoriesId: g.CategoriesId,
		IsValidated:  g.IsValidated,
	}
}

// UpdateGenre contains updates of an existing genre.
type UpdateGenre struct {
	Name         *string
	CategoriesId []string
	IsValidated  *bool
}

func (g UpdateGenre) update(genre Genre) Genre {
	if g.Name != nil {
		genre.Name = *g.Name
	}
	if g.CategoriesId != nil {
		genre.CategoriesId = g.CategoriesId
	}
	if g.IsValidated != nil {
		genre.IsValidated = *g.IsValidated
	}
	return genre
}

type service struct {
	idGenerator id_generator.IdGenerator
	repo        Repository
}

// NewService returns a new Service with all of the expected middlewares wired in.
func NewService(idGenerator id_generator.IdGenerator, r Repository, logger log.Logger) Service {
	var svc Service
	{
		svc = service{idGenerator: idGenerator, repo: r}
		svc = NewValidationMiddleware()(svc)
		svc = NewLoggingMiddleware(logger)(svc)
	}
	return svc
}

func (svc service) Create(ctx context.Context, newGenre NewGenre) (Genre, error) {
	id, err := svc.idGenerator.Generate()
	if err != nil {
		return Genre{}, err
	}
	genre := newGenre.toGenre(id)
	err = svc.repo.Store(ctx, genre)
	if err != nil {
		return Genre{}, err
	}
	return genre, nil
}

func (svc service) Destroy(ctx context.Context, id string) error {
	return svc.repo.DeleteOne(ctx, id)
}

func (svc service) List(ctx context.Context) ([]Genre, error) {
	genres, err := svc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return genres, nil
}

func (svc service) Show(ctx context.Context, id string) (Genre, error) {
	genre, err := svc.repo.GetOne(ctx, id)
	if err != nil {
		return Genre{}, err
	}
	return genre, nil
}

func (svc service) Update(ctx context.Context, id string, updateGenre UpdateGenre) error {
	err := svc.repo.UpdateOne(ctx, id, updateGenre)
	if err != nil {
		return err
	}
	return nil
}

func (g NewGenre) isEmpty() bool {
	return g.compare(NewGenre{})
}

func (g NewGenre) compare(b NewGenre) bool {
	if &g == &b {
		return true
	}
	if g.Name != b.Name {
		return false
	}
	if len(g.CategoriesId) != len(b.CategoriesId) {
		return false
	}
	for i, v := range g.CategoriesId {
		if b.CategoriesId[i] != v {
			return false
		}
	}
	return true
}

func (g UpdateGenre) isEmpty() bool {
	return g.compare(UpdateGenre{})
}

func (g UpdateGenre) compare(b UpdateGenre) bool {
	if &g == &b {
		return true
	}
	if g.Name != b.Name {
		return false
	}
	if len(g.CategoriesId) != len(b.CategoriesId) {
		return false
	}
	for i, v := range g.CategoriesId {
		if b.CategoriesId[i] != v {
			return false
		}
	}
	return true
}
