package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/selmison/code-micro-videos/pkg/crud/domain"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (s service) CreateGenre(ctx domain.Context, fields domain.GenreValidatable) error {
	genre, err := domain.NewGenre(fields)
	if err != nil {
		return fmt.Errorf("error CreateGenre(): %w", err)
	}
	return s.r.CreateGenre(ctx, *genre)
}

func (s service) FetchGenre(ctx domain.Context, name string) (domain.Genre, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	if len(name) == 0 {
		return domain.Genre{}, fmt.Errorf("'name' %w", logger.ErrIsRequired)
	}
	genre, err := s.r.FetchGenre(ctx, name)
	if err == sql.ErrNoRows {
		return domain.Genre{}, fmt.Errorf("%s: %w", name, logger.ErrNotFound)
	} else if err != nil {
		return domain.Genre{}, fmt.Errorf("%s: %w", name, logger.ErrInternalApplication)
	}
	return genre, nil
}

func (s service) GetGenres(ctx domain.Context, limit int) ([]domain.Genre, error) {
	if limit < 0 {
		return nil, logger.ErrInvalidedLimit
	}
	return s.r.GetGenres(ctx, limit)
}

func (s service) RemoveGenre(ctx domain.Context, name string) error {
	name = strings.ToLower(strings.TrimSpace(name))
	if len(name) == 0 {
		return fmt.Errorf("'name' %w", logger.ErrIsRequired)
	}
	if err := s.r.RemoveGenre(ctx, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", name, logger.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s service) UpdateGenre(ctx domain.Context, name string, fields domain.GenreValidatable) error {
	name = strings.ToLower(strings.TrimSpace(name))
	genre, err := domain.NewGenre(fields)
	if err != nil {
		return fmt.Errorf("error UpdateGenre(): %w", err)
	}
	if err := s.r.UpdateGenre(ctx, name, *genre); err != nil {
		return err
	}
	return nil
}
