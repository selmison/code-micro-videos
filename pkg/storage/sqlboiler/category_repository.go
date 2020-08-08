package sqlboiler

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (r Repository) UpdateCategory(name string, c crud.CategoryDTO) error {
	category, err := r.FetchCategory(name)
	if err != nil {
		return err
	}
	nameDTO := strings.ToLower(strings.TrimSpace(c.Name))
	category.Name = nameDTO
	DescriptionDTO := strings.TrimSpace(c.Description)
	category.Description = null.String{String: DescriptionDTO, Valid: true}
	_, err = category.UpdateG(r.ctx, boil.Infer())
	if err != nil {
		return fmt.Errorf("%s %w", nameDTO, logger.ErrAlreadyExists)
	}
	return nil
}

func (r Repository) AddCategory(categoryDTO crud.CategoryDTO) error {
	newCat := models.Category{
		ID:          uuid.New().String(),
		Name:        strings.ToLower(strings.TrimSpace(categoryDTO.Name)),
		Description: null.String{String: categoryDTO.Description, Valid: true},
	}
	err := newCat.InsertG(r.ctx, boil.Infer())
	if err != nil {
		var e *pq.Error
		if errors.As(err, &e) {
			if e.Code.Name() == "unique_violation" {
				return fmt.Errorf("name '%s' %w", categoryDTO.Name, logger.ErrAlreadyExists)
			} else {
				return fmt.Errorf("%s: %w", "method Repository.AddCategory(categoryDTO)", err)
			}
		}
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
	name = strings.ToLower(strings.TrimSpace(name))
	categorySlice, err := models.Categories(Where("is_validated=?", true), models.CategoryWhere.Name.EQ(name)).AllG(r.ctx)
	if err != nil {
		return models.Category{}, err
	}
	if len(categorySlice) == 0 {
		return models.Category{}, sql.ErrNoRows
	}
	return *categorySlice[0], nil
}
