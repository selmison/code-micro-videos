package sqlboiler

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (r Repository) UpdateGenre(name string, genreDTO crud.GenreDTO) error {
	genre, err := r.FetchGenre(name)
	if err != nil {
		return err
	}
	nameDTO := strings.ToLower(strings.TrimSpace(genreDTO.Name))
	genre.Name = nameDTO
	_, err = genre.UpdateG(r.ctx, boil.Infer())
	if err != nil {
		return fmt.Errorf("%s %w", nameDTO, logger.ErrAlreadyExists)
	}
	return nil
}

func (r Repository) AddGenre(genreDTO crud.GenreDTO) error {
	genre := models.Genre{
		ID:   uuid.New().String(),
		Name: strings.ToLower(strings.TrimSpace(genreDTO.Name)),
	}
	err := genre.InsertG(r.ctx, boil.Infer())
	if err != nil {
		var e *pq.Error
		if errors.As(err, &e) {
			if e.Code.Name() == "unique_violation" {
				return fmt.Errorf("name '%s' %w", genreDTO.Name, logger.ErrAlreadyExists)
			} else {
				return fmt.Errorf("%s: %w", "method Repository.AddGenre(genreDTO)", err)
			}
		}
	}
	return nil
}

func (r Repository) RemoveGenre(name string) error {
	c, err := r.FetchGenre(name)
	if err != nil {
		return err
	}
	c.IsValidated = false
	_, err = c.UpdateG(r.ctx, boil.Infer())
	return err
}

func (r Repository) GetGenres(limit int) (models.GenreSlice, error) {
	if limit <= 0 {
		return nil, nil
	}
	genres, err := models.Genres(Where("is_validated=?", true), Limit(limit)).AllG(r.ctx)
	if err != nil {
		return nil, err
	}
	return genres, nil
}

func (r Repository) FetchGenre(name string) (models.Genre, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	genreSlice, err := models.Genres(Where("is_validated=?", true), models.GenreWhere.Name.EQ(name)).AllG(r.ctx)
	if err != nil {
		return models.Genre{}, err
	}
	if len(genreSlice) == 0 {
		return models.Genre{}, sql.ErrNoRows
	}
	return *genreSlice[0], nil
}
