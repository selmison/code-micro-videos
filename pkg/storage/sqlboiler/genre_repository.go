package sqlboiler

import (
	"database/sql"
	"errors"
	"fmt"

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
	genre.Name = genreDTO.Name
	tx, err := boil.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}
	_, err = genre.Update(r.ctx, tx, boil.Infer())
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("%s %w", genreDTO.Name, logger.ErrAlreadyExists)
	}
	if err := r.setCategoriesInGenre(genreDTO.Categories, genre, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r Repository) AddGenre(genreDTO crud.GenreDTO) error {
	genre := models.Genre{
		ID:   uuid.New().String(),
		Name: genreDTO.Name,
	}
	tx, err := boil.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}
	err = genre.Insert(r.ctx, tx, boil.Infer())
	if err != nil {
		var e *pq.Error
		if err := tx.Rollback(); err != nil {
			return err
		}
		if errors.As(err, &e) {
			if e.Code.Name() == "unique_violation" {
				return fmt.Errorf("name '%s' %w", genreDTO.Name, logger.ErrAlreadyExists)
			} else {
				return fmt.Errorf("%s: %w", "method Repository.AddGenre(genreDTO)", err)
			}
		}
	}
	if err := r.setCategoriesInGenre(genreDTO.Categories, genre, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
func (r Repository) setCategoriesInGenre(categories []crud.CategoryDTO, genre models.Genre, tx *sql.Tx) error {
	if categories == nil || len(categories) == 0 {
		return nil
	}
	clause := "name=?"
	categoryNames := make([]interface{}, len(categories))
	for i, category := range categories {
		if i > 0 {
			clause = fmt.Sprintf("name=? OR %s", clause)
		}
		categoryNames[i] = category.Name
	}
	categorySlice, err := models.Categories(
		Where(clause, categoryNames...),
	).AllG(r.ctx)
	if err != nil {
		return err
	}
	if len(categorySlice) == 0 {
		return fmt.Errorf("none category is %w", logger.ErrNotFound)
	}
	if err := genre.SetCategories(r.ctx, tx, false, categorySlice...); err != nil {
		return fmt.Errorf("insert a new slice of categories and assign them to the category: %s", err)
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
	genreSlice, err := models.Genres(Where("is_validated=?", true), models.GenreWhere.Name.EQ(name)).AllG(r.ctx)
	if err != nil {
		return models.Genre{}, err
	}
	if len(genreSlice) == 0 {
		return models.Genre{}, sql.ErrNoRows
	}
	return *genreSlice[0], nil
}
