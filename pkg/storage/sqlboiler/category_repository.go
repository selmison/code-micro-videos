package sqlboiler

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud/domain"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (r Repository) CreateCategory(ctx domain.Context, category domain.Category) error {
	categoryValidatable := category.MapToCategoryValidatable()
	categorySqlBoiler := models.Category{
		ID:          categoryValidatable.Id,
		Name:        categoryValidatable.Name,
		Description: null.String{String: categoryValidatable.Description, Valid: true},
	}
	tx, err := boil.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}
	err = categorySqlBoiler.Insert(r.ctx, tx, boil.Infer())
	if err != nil {
		var e *pq.Error
		if err := tx.Rollback(); err != nil {
			return err
		}
		if errors.As(err, &e) {
			if e.Code.Name() == "unique_violation" {
				return fmt.Errorf("name '%s' %w", categoryValidatable.Name, logger.ErrAlreadyExists)
			} else {
				return fmt.Errorf("%s: %w", "method Repository.CreateCategory(category)", err)
			}
		}
	}
	if err := r.setGenresInCategory(ctx, categoryValidatable.Genres, categorySqlBoiler, tx); err != nil {
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

func (r Repository) FetchCategory(ctx domain.Context, name string) (domain.Category, error) {
	categorySqlBoiler, err := r.fetchCategorySqlBoiler(ctx, name)
	if err != nil {
		return domain.Category{}, err
	}
	categoryDomain, err := mapToDomainCategory(*categorySqlBoiler)
	if err != nil {
		return domain.Category{}, err
	}
	return *categoryDomain, nil
}

func (r Repository) fetchCategorySqlBoiler(ctx domain.Context, name string) (*models.Category, error) {
	categorySlice, err := models.Categories(Where("is_validated=?", true), models.CategoryWhere.Name.EQ(name)).AllG(ctx)
	if err != nil {
		return nil, err
	}
	if len(categorySlice) == 0 {
		return nil, sql.ErrNoRows
	}
	return categorySlice[0], nil
}

func (r Repository) GetCategories(ctx domain.Context, limit int) ([]domain.Category, error) {
	if limit <= 0 {
		return nil, nil
	}
	categorySqlBoiler, err := models.Categories(Where("is_validated=?", true), Limit(limit)).AllG(ctx)
	if err != nil {
		return nil, err
	}
	categories := make([]domain.Category, len(categorySqlBoiler))
	for i, categorySqlBoiler := range categorySqlBoiler {
		categoryDomain, err := mapToDomainCategory(*categorySqlBoiler)
		if err != nil {
			return nil, err
		}
		categories[i] = *categoryDomain
	}
	return categories, nil
}

func (r Repository) GetCategoriesByIds(ctx domain.Context, ids []interface{}) ([]domain.Category, error) {
	if ids == nil || len(ids) == 0 {
		return nil, fmt.Errorf("none category is %w", logger.ErrNotFound)
	}
	categorySlice, err := models.Categories(
		WhereIn("is_validated=?", true),
		AndIn("id in ?", ids...),
	).AllG(ctx)
	if err != nil {
		return nil, fmt.Errorf("error GetCategoriesByIds(): %w\n", err)
	}
	categories := make([]domain.Category, len(categorySlice))
	for i, categorySqlBoiler := range categorySlice {
		categoryDomain, err := mapToDomainCategory(*categorySqlBoiler)
		if err != nil {
			return nil, err
		}
		categories[i] = *categoryDomain
	}
	models.Categories(
		Where("is_validated=?", true),
		AndIn("id in ?", ids...),
	)
	return categories, nil
}

func mapToDomainCategory(category models.Category) (*domain.Category, error) {
	var genres []domain.Genre
	if category.R == nil {
		genres = nil
	} else {
		genres = mapToGenresValidatable(category.R.Genres)
	}
	categoryValidatable := domain.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description.String,
		Genres:      genres,
	}
	return domain.NewCategory(categoryValidatable)
}

func mapToCategoriesValidatable(categoriesSqlBoiler models.CategorySlice) []domain.Category {
	var categoriesValidatable []domain.Category
	if categoriesSqlBoiler == nil || len(categoriesSqlBoiler) == 0 {
		return nil
	} else {
		categoriesValidatable = make([]domain.Category, len(categoriesSqlBoiler))
		for j, category := range categoriesSqlBoiler {
			genres := mapToGenresValidatable(category.R.Genres)
			categoryValidatable := domain.Category{
				Id:          category.ID,
				Name:        category.Name,
				Description: category.Description.String,
				Genres:      genres,
			}
			categoriesValidatable[j] = categoryValidatable
		}
	}
	return categoriesValidatable
}

func (r Repository) RemoveCategoryByName(ctx domain.Context, name string) error {
	c, err := r.fetchCategorySqlBoiler(ctx, name)
	if err != nil {
		return err
	}
	c.IsValidated = false
	_, err = c.UpdateG(r.ctx, boil.Infer())
	return err
}

func (r Repository) setGenresInCategory(ctx domain.Context, gs []domain.Genre, category models.Category, tx *sql.Tx) error {
	if gs == nil || len(gs) == 0 {
		return nil
	}
	var genreIds []interface{}
	for _, genre := range gs {
		genreIds = append(genreIds, genre.Id)
	}
	genres, err := r.GetGenresByIds(ctx, genreIds)
	if err != nil {
		return err
	}
	clause := "name=?"
	genreNames := make([]interface{}, len(genres))
	for i, genre := range genres {
		if i > 0 {
			clause = fmt.Sprintf("name=? OR %s", clause)
		}
		genreNames[i] = genre.Name()
	}
	genreSlice, err := models.Genres(
		Where(clause, genreNames...),
	).AllG(ctx)
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

func (r Repository) UpdateCategory(ctx domain.Context, name string, category domain.Category) error {
	categoryValidatable := category.MapToCategoryValidatable()
	categorySqlBoiler, err := r.fetchCategorySqlBoiler(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", name, logger.ErrNotFound)
		}
		return err
	}
	categorySqlBoiler.Name = categoryValidatable.Name
	tx, err := boil.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}
	_, err = categorySqlBoiler.Update(r.ctx, tx, boil.Infer())
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("%s %w", categoryValidatable.Name, logger.ErrAlreadyExists)
	}
	if err := r.setGenresInCategory(ctx, categoryValidatable.Genres, *categorySqlBoiler, tx); err != nil {
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
