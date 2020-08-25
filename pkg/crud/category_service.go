//go:generate mockgen -destination=./mock/service.go -package=mock . Repository,Service

package crud

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

type service struct {
	r Repository
}

func (s service) RemoveCategory(name string) error {
	name = strings.ToLower(strings.TrimSpace(name))
	if len(strings.TrimSpace(name)) == 0 {
		return fmt.Errorf("'name' %w", logger.ErrIsRequired)
	}
	if err := s.r.RemoveCategory(name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", name, logger.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s service) UpdateCategory(name string, c CategoryDTO) error {
	name = strings.ToLower(strings.TrimSpace(name))
	if len(name) == 0 {
		return fmt.Errorf("'name' %w", logger.ErrIsRequired)
	}
	c.Name = strings.ToLower(strings.TrimSpace(c.Name))
	if err := c.Validate(); err != nil {
		return err
	}
	if err := s.r.UpdateCategory(name, c); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", name, logger.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s service) AddCategory(c CategoryDTO) error {
	c.Name = strings.ToLower(strings.TrimSpace(c.Name))
	if err := c.Validate(); err != nil {
		return err
	}
	return s.r.AddCategory(c)
}
func (s service) GetCategories(limit int) (models.CategorySlice, error) {
	if limit < 0 {
		return nil, logger.ErrInvalidedLimit
	}
	return s.r.GetCategories(limit)
}

func (s service) FetchCategory(name string) (models.Category, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	c, err := s.r.FetchCategory(name)
	if err == sql.ErrNoRows {
		return models.Category{}, fmt.Errorf("%s: %w", name, logger.ErrNotFound)
	}
	return c, nil
}
