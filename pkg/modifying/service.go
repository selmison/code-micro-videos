package modifying

import (
	"fmt"
	"strings"
)

var (
	ErrIsRequired    = fmt.Errorf("is required")
	ErrAlreadyExists = fmt.Errorf("already exists")
)

type Repository interface {
	Service
}

type Service interface {
	AddCategory(c CategoryDTO) error
	RemoveCategory(name string) error
	UpdateCategory(name string, c CategoryDTO) error
}

type service struct {
	r Repository
}

func (s service) RemoveCategory(name string) error {
	if len(strings.TrimSpace(name)) == 0 {
		return fmt.Errorf("'name' %w", ErrIsRequired)

	}
	return s.r.RemoveCategory(name)
}

func (s service) UpdateCategory(name string, c CategoryDTO) error {
	if len(strings.TrimSpace(name)) == 0 {
		return fmt.Errorf("'name' %w", ErrIsRequired)
	}
	if err := c.Validate(); err != nil {
		return err
	}
	return s.r.UpdateCategory(name, c)
}

func (s service) AddCategory(c CategoryDTO) error {
	c.Name = strings.ToLower(strings.TrimSpace(c.Name))
	if err := c.Validate(); err != nil {
		return err
	}
	return s.r.AddCategory(c)
}

// NewService creates a listing service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}
