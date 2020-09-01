package sqlboiler

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (r Repository) UpdateCategory(name string, categoryDTO crud.CategoryDTO) error {
	category, err := r.FetchCategory(name)
	if err != nil {
		return err
	}
	category.Name = categoryDTO.Name
	category.Description = null.String{String: categoryDTO.Description, Valid: true}
	tx, err := boil.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}
	_, err = category.Update(r.ctx, tx, boil.Infer())
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("%s %w", categoryDTO.Name, logger.ErrAlreadyExists)
	}
	if err := r.setGenresInCategory(categoryDTO.Genres, category, tx); err != nil {
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

func (r Repository) AddCategory(categoryDTO crud.CategoryDTO) error {
	category := models.Category{
		ID:          uuid.New().String(),
		Name:        categoryDTO.Name,
		Description: null.String{String: categoryDTO.Description, Valid: true},
	}
	tx, err := boil.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}
	err = category.Insert(r.ctx, tx, boil.Infer())
	if err != nil {
		var e *pq.Error
		if err := tx.Rollback(); err != nil {
			return err
		}
		if errors.As(err, &e) {
			if e.Code.Name() == "unique_violation" {
				return fmt.Errorf("name '%s' %w", categoryDTO.Name, logger.ErrAlreadyExists)
			} else {
				return fmt.Errorf("%s: %w", "method Repository.AddCategory(categoryDTO)", err)
			}
		}
	}
	if err := r.setGenresInCategory(categoryDTO.Genres, category, tx); err != nil {
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

func (r Repository) setGenresInCategory(genres []crud.GenreDTO, category models.Category, tx *sql.Tx) error {
	if genres == nil || len(genres) == 0 {
		return nil
	}
	clause := "name=?"
	genreNames := make([]interface{}, len(genres))
	for i, genre := range genres {
		if i > 0 {
			clause = fmt.Sprintf("name=? OR %s", clause)
		}
		genreNames[i] = genre.Name
	}
	genreSlice, err := models.Genres(
		Where(clause, genreNames...),
	).AllG(r.ctx)
	if err != nil {
		return err
	}
	if len(genreSlice) == 0 {
		return fmt.Errorf("none genre is %w", logger.ErrNotFound)
	}
	if err := category.SetGenres(r.ctx, tx, false, genreSlice...); err != nil {
		return fmt.Errorf("insert a new slice of genres and assign them to the category: %s", err)
	}
	return nil
}

func (r Repository) RemoveCategory(name string) error {
	c, err := r.FetchCategory(name)
	if err != nil {
		return err
	}
	c.IsValidated = false
	_, err = c.UpdateG(r.ctx, boil.Infer())
	return err
}

func (r Repository) GetCategories(limit int) (models.CategorySlice, error) {
	if limit <= 0 {
		return nil, nil
	}
	categories, err := models.Categories(Where("is_validated=?", true), Limit(limit)).AllG(r.ctx)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r Repository) FetchCategory(name string) (models.Category, error) {
	categorySlice, err := models.Categories(Where("is_validated=?", true), models.CategoryWhere.Name.EQ(name)).AllG(r.ctx)
	if err != nil {
		return models.Category{}, err
	}
	if len(categorySlice) == 0 {
		return models.Category{}, sql.ErrNoRows
	}
	return *categorySlice[0], nil
}
