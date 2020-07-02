package listing

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/selmison/code-micro-videos/models"
)

var (
	ErrInvalidedLimit = fmt.Errorf("limit should be positive number")
	ErrNotFound       = errors.New("not found")
	validate          *validator.Validate
)

type Repository interface {
	Service
}

type Service interface {
	GetCategories(limit int) (models.CategorySlice, error)
	FetchCategory(name string) (models.Category, error)
}

type service struct {
	r Repository
}

func (s service) GetCategories(limit int) (models.CategorySlice, error) {
	if limit < 0 {
		return nil, ErrInvalidedLimit
	}
	return s.r.GetCategories(limit)
}

func (s service) FetchCategory(name string) (models.Category, error) {
	c, error := s.r.FetchCategory(name)
	if error == sql.ErrNoRows {
		return models.Category{}, fmt.Errorf("%s: %w", name, ErrNotFound)
	}
	return c, nil
}

// NewService creates a listing service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

func init() {
	validate = validator.New()
}
