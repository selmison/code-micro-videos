package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/selmison/code-micro-videos/pkg/crud/domain"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (s service) CreateCategory(ctx domain.Context, fields domain.CategoryValidatable) error {
	category, err := domain.NewCategory(fields)
	if err != nil {
		return fmt.Errorf("error CreateCategory(): %w", err)
	}
	return s.r.CreateCategory(ctx, *category)
}

func (s service) FetchCategory(ctx domain.Context, name string) (domain.Category, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	if len(name) == 0 {
		return domain.Category{}, fmt.Errorf("'name' %w", logger.ErrIsRequired)
	}
	category, err := s.r.FetchCategory(ctx, name)
	if err == sql.ErrNoRows {
		return domain.Category{}, fmt.Errorf("%s: %w", name, logger.ErrNotFound)
	} else if err != nil {
		return domain.Category{}, fmt.Errorf("%s: %w", name, logger.ErrInternalApplication)
	}
	return category, nil
}

func (s service) GetCategories(ctx domain.Context, limit int) ([]domain.Category, error) {
	if limit < 0 {
		return nil, logger.ErrInvalidedLimit
	}
	return s.r.GetCategories(ctx, limit)
}

func (s service) RemoveCategory(ctx domain.Context, name string) error {
	name = strings.ToLower(strings.TrimSpace(name))
	if len(name) == 0 {
		return fmt.Errorf("'name' %w", logger.ErrIsRequired)
	}
	if err := s.r.RemoveCategory(ctx, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", name, logger.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s service) UpdateCategory(ctx domain.Context, name string, fields domain.CategoryValidatable) error {
	name = strings.ToLower(strings.TrimSpace(name))
	category, err := domain.NewCategory(fields)
	if err != nil {
		return fmt.Errorf("error UpdateCategory(): %w", err)
	}
	if err := s.r.UpdateCategory(ctx, name, *category); err != nil {
		return err
	}
	return nil
}
