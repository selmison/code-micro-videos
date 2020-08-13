package testdata

import (
	"github.com/bxcodec/faker/v3"
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
	FakeGenres = []models.Genre{
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
	FakeCastMembers = []models.CastMember{
		{
			ID:   uuid.New().String(),
			Name: faker.Name(),
			Type: int16(crud.Actor),
		},
		{
			ID:   uuid.New().String(),
			Name: faker.Name(),
			Type: int16(crud.Director),
		},
		{
			ID:   uuid.New().String(),
			Name: faker.Name(),
			Type: int16(crud.Actor),
		},
		{
			ID:   uuid.New().String(),
			Name: faker.Name(),
			Type: int16(crud.Actor),
		},
		{
			ID:   uuid.New().String(),
			Name: faker.Name(),
			Type: int16(crud.Director),
		},
		{
			ID:   uuid.New().String(),
			Name: faker.Name(),
			Type: int16(crud.Actor),
		},
	}
	FakeCategoriesDTO  []crud.CategoryDTO
	FakeGenresDTO      []crud.GenreDTO
	FakeCastMembersDTO []crud.CastMemberDTO
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
	FakeCastMembersDTO = make([]crud.CastMemberDTO, len(FakeCastMembers))
	for i, castMember := range FakeCastMembers {
		FakeCastMembersDTO[i] = crud.CastMemberDTO{
			Name: castMember.Name,
			Type: crud.CastMemberType(castMember.Type),
		}
	}
}
