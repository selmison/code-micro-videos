package testdata

import (
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud"
)

var (
	FakeCategories = []models.Category{
		{
			ID:          uuid.New().String(),
			Name:        "action",
			Description: null.String{String: "action films", Valid: true},
		},
		{
			ID:   uuid.New().String(),
			Name: "animation",
		},
		{
			ID:          uuid.New().String(),
			Name:        "science fiction",
			Description: null.String{String: "science fiction films", Valid: true},
		},
		{
			ID:   uuid.New().String(),
			Name: "violent",
		},
		{
			ID:          uuid.New().String(),
			Name:        "drama",
			Description: null.String{String: "drama films", Valid: true},
		},
		{
			ID:   uuid.New().String(),
			Name: "romance",
		},
	}
	FakeCategoriesDTO []crud.CategoryDTO
	FakeGenres        = []models.Genre{
		{
			ID:   uuid.New().String(),
			Name: "action",
		},
		{
			ID:   uuid.New().String(),
			Name: "animation",
		},
		{
			ID:   uuid.New().String(),
			Name: "science fiction",
		},
		{
			ID:   uuid.New().String(),
			Name: "violent",
		},
		{
			ID:   uuid.New().String(),
			Name: "drama",
		},
		{
			ID:   uuid.New().String(),
			Name: "romance",
		},
	}
	FakeGenresDTO []crud.GenreDTO
)

func init() {
	FakeCategoriesDTO = make([]crud.CategoryDTO, len(FakeCategories))
	for i, category := range FakeCategories {
		FakeCategoriesDTO[i] = crud.CategoryDTO{
			Name:        category.Name,
			Description: category.Description.String,
		}
	}
	FakeGenresDTO = make([]crud.GenreDTO, len(FakeGenres))
	for i, user := range FakeGenres {
		FakeGenresDTO[i] = crud.GenreDTO{
			Name: user.Name,
		}
	}
}
