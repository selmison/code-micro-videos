package crud

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (s service) RemoveGenre(name string) error {
	if len(strings.TrimSpace(name)) == 0 {
		return fmt.Errorf("'name' %w", logger.ErrIsRequired)
	}
	if err := s.r.RemoveGenre(name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", name, logger.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s service) UpdateGenre(name string, genreDTO GenreDTO) error {
	if len(strings.TrimSpace(name)) == 0 {
		return fmt.Errorf("'name' %w", logger.ErrIsRequired)
	}
	if err := genreDTO.Validate(); err != nil {
		return err
	}
	if err := s.r.UpdateGenre(name, genreDTO); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", name, logger.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s service) AddGenre(genreDTO GenreDTO) error {
	genreDTO.Name = strings.ToLower(strings.TrimSpace(genreDTO.Name))
	if err := genreDTO.Validate(); err != nil {
		return err
	}
	return s.r.AddGenre(genreDTO)
}
func (s service) GetGenres(limit int) (models.GenreSlice, error) {
	if limit < 0 {
		return nil, logger.ErrInvalidedLimit
	}
	return s.r.GetGenres(limit)
}

func (s service) FetchGenre(name string) (models.Genre, error) {
	c, err := s.r.FetchGenre(name)
	if err == sql.ErrNoRows {
		return models.Genre{}, fmt.Errorf("%s: %w", name, logger.ErrNotFound)
	} else if err != nil {
		return models.Genre{}, fmt.Errorf("%s: %w", name, logger.ErrInternalApplication)
	}

	return c, nil
}
