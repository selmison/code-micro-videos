package sqlboiler

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/modifying"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/boil"
	. "github.com/volatiletech/sqlboiler/queries/qm"
	"strings"
)

type Repository struct {
	ctx context.Context
}

func (r Repository) UpdateCategory(name string, c modifying.CategoryDTO) error {
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
		return fmt.Errorf("%s %w", nameDTO, modifying.ErrAlreadyExists)
	}
	return nil
}

func (r Repository) AddCategory(c modifying.CategoryDTO) error {
	newCat := models.Category{
		ID:          uuid.New().String(),
		Name:        strings.ToLower(strings.TrimSpace(c.Name)),
		Description: null.String{String: c.Description, Valid: true},
	}
	err := newCat.InsertG(r.ctx, boil.Infer())
	if err != nil {
		return fmt.Errorf("%s: %w", c.Name, modifying.ErrAlreadyExists)
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

func NewRepository(ctx context.Context, db *sql.DB) Repository {
	boil.SetDB(db)
	return Repository{ctx}
}
