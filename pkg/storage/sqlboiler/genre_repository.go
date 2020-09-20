package sqlboiler

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud/domain"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (r Repository) CreateGenre(ctx domain.Context, genre domain.Genre) error {
	genreValidatable := genre.MapToGenreValidatable()
	genreSqlBoiler := models.Genre{
		ID:   genreValidatable.Id,
		Name: genreValidatable.Name,
	}
	tx, err := boil.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}
	err = genreSqlBoiler.Insert(r.ctx, tx, boil.Infer())
	if err != nil {
		var e *pq.Error
		if err := tx.Rollback(); err != nil {
			return err
		}
		if errors.As(err, &e) {
			if e.Code.Name() == "unique_violation" {
				return fmt.Errorf("name '%s' %w", genreValidatable.Name, logger.ErrAlreadyExists)
			} else {
				return fmt.Errorf("%s: %w", "method Repository.CreateGenre(genre)", err)
			}
		}
	}
	if err := r.setCategoriesInGenre(ctx, genreValidatable.Categories, genreSqlBoiler, tx); err != nil {
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

func (r Repository) FetchGenre(ctx domain.Context, name string) (domain.Genre, error) {
	genreSqlBoiler, err := r.fetchGenreSqlBoiler(ctx, name)
	if err != nil {
		return domain.Genre{}, err
	}
	genreDomain, err := mapToDomainGenre(*genreSqlBoiler)
	if err != nil {
		return domain.Genre{}, err
	}
	return *genreDomain, nil
}

func (r Repository) fetchGenreSqlBoiler(ctx domain.Context, name string) (*models.Genre, error) {
	genreSlice, err := models.Genres(Where("is_validated=?", true), models.GenreWhere.Name.EQ(name)).AllG(ctx)
	if err != nil {
		return nil, err
	}
	if len(genreSlice) == 0 {
		return nil, sql.ErrNoRows
	}
	return genreSlice[0], nil
}

func (r Repository) GetGenres(ctx domain.Context, limit int) ([]domain.Genre, error) {
	if limit <= 0 {
		return nil, nil
	}
	genreSqlBoiler, err := models.Genres(Where("is_validated=?", true), Limit(limit)).AllG(ctx)
	if err != nil {
		return nil, err
	}
	genres := make([]domain.Genre, len(genreSqlBoiler))
	for i, genreSqlBoiler := range genreSqlBoiler {
		genreDomain, err := mapToDomainGenre(*genreSqlBoiler)
		if err != nil {
			return nil, err
		}
		genres[i] = *genreDomain
	}
	return genres, nil
}

func (r Repository) GetGenresByIds(ctx domain.Context, ids []interface{}) ([]domain.Genre, error) {
	if ids == nil || len(ids) == 0 {
		return nil, fmt.Errorf("none genre is %w", logger.ErrNotFound)
	}
	genreSlice, err := models.Genres(
		WhereIn("is_validated=?", true),
		AndIn("id in ?", ids...),
	).AllG(ctx)
	if err != nil {
		return nil, fmt.Errorf("error GetGenresByIds(): %w\n", err)
	}
	genres := make([]domain.Genre, len(genreSlice))
	for i, genreSqlBoiler := range genreSlice {
		genreDomain, err := mapToDomainGenre(*genreSqlBoiler)
		if err != nil {
			return nil, err
		}
		genres[i] = *genreDomain
	}
	models.Genres(
		Where("is_validated=?", true),
		AndIn("id in ?", ids...),
	)
	return genres, nil
}

func mapToDomainGenre(genre models.Genre) (*domain.Genre, error) {
	var categories []domain.Category
	if genre.R == nil {
		categories = nil
	} else {
		categories = mapToCategoriesValidatable(genre.R.Categories)
	}
	genreValidatable := domain.Genre{
		Id:         genre.ID,
		Name:       genre.Name,
		Categories: categories,
	}
	return domain.NewGenre(genreValidatable)
}

func mapToGenresValidatable(genresSqlBoiler models.GenreSlice) []domain.Genre {
	var genreValidatable []domain.Genre
	if genresSqlBoiler == nil || len(genresSqlBoiler) == 0 {
		return nil
	} else {
		genreValidatable = make([]domain.Genre, len(genresSqlBoiler))
		for j, genre := range genresSqlBoiler {
			categories := mapToCategoriesValidatable(genre.R.Categories)
			categoryValidatable := domain.Genre{
				Id:         genre.ID,
				Name:       genre.Name,
				Categories: categories,
			}
			genreValidatable[j] = categoryValidatable
		}
	}
	return genreValidatable
}

func (r Repository) RemoveGenreByName(ctx domain.Context, name string) error {
	c, err := r.fetchGenreSqlBoiler(ctx, name)
	if err != nil {
		return err
	}
	c.IsValidated = false
	_, err = c.UpdateG(r.ctx, boil.Infer())
	return err
}

func (r Repository) setCategoriesInGenre(ctx domain.Context, cs []domain.Category, genre models.Genre, tx *sql.Tx) error {
	if cs == nil || len(cs) == 0 {
		return nil
	}
	var categoryIds []interface{}
	for _, category := range cs {
		categoryIds = append(categoryIds, category.Id)
	}
	categories, err := r.GetCategoriesByIds(ctx, categoryIds)
	if err != nil {
		return err
	}
	clause := "name=?"
	categoryNames := make([]interface{}, len(categories))
	for i, category := range categories {
		if i > 0 {
			clause = fmt.Sprintf("name=? OR %s", clause)
		}
		categoryNames[i] = category.Name()
	}
	categorySlice, err := models.Categories(
		Where(clause, categoryNames...),
	).AllG(ctx)
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

func (r Repository) UpdateGenre(ctx domain.Context, name string, genre domain.Genre) error {
	genreValidatable := genre.MapToGenreValidatable()
	genreSqlBoiler, err := r.fetchGenreSqlBoiler(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", name, logger.ErrNotFound)
		}
		return err
	}
	genreSqlBoiler.Name = genreValidatable.Name
	tx, err := boil.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}
	_, err = genreSqlBoiler.Update(r.ctx, tx, boil.Infer())
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("%s %w", genreValidatable.Name, logger.ErrAlreadyExists)
	}
	if err := r.setCategoriesInGenre(ctx, genreValidatable.Categories, *genreSqlBoiler, tx); err != nil {
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
