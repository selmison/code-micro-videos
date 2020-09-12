package domain

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

type GenreValidatable struct {
	Id         string `validate:"not_blank"`
	Name       string `validate:"not_blank"`
	Categories []CategoryValidatable
	validated  bool
}

type Genre struct {
	id         string
	name       string
	categories []Category
}

func (g *Genre) Name() string {
	return g.name
}

func NewGenre(fields GenreValidatable) (*Genre, error) {
	return fields.mapToGenre()
}

func (g *Genre) MapToGenreValidatable() *GenreValidatable {
	var categories []CategoryValidatable
	if g.categories != nil && len(g.categories) > 0 {
		categories = make([]CategoryValidatable, len(g.categories))
		for i, category := range g.categories {
			categories[i] = *category.MapToCategoryValidatable()
		}

	}
	return &GenreValidatable{
		Id:         g.id,
		Name:       g.name,
		Categories: categories,
		validated:  true,
	}
}

func mapToGenreValidatables(genres []Genre) []GenreValidatable {
	var gValidatables []GenreValidatable
	if len(genres) > 0 {
		gValidatables = make([]GenreValidatable, len(genres))
		for i, genre := range genres {
			gValidatables[i] = *genre.MapToGenreValidatable()
		}
	}
	return gValidatables
}

func (g *GenreValidatable) IsValidated() bool {
	return g.validated
}

func (g *GenreValidatable) mapToGenre() (*Genre, error) {
	if err := g.Validate(); err != nil {
		return nil, err
	}
	var categories []Category
	if len(g.Categories) > 0 {
		categories = make([]Category, len(g.Categories))
		for i, categoryValidatable := range g.Categories {
			category, err := categoryValidatable.mapToCategory()
			if err != nil {
				return nil, err
			}
			categories[i] = *category
		}
	}
	return &Genre{
		id:         g.Id,
		name:       g.Name,
		categories: categories,
	}, nil
}

func (g *GenreValidatable) Validate() error {
	g.validated = false
	if strings.ToLower(g.Name) != g.Name {
		return fmt.Errorf("lowercase 'Name' field %v", logger.ErrIsRequired)
	}
	err := validate.Struct(g)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		return fmt.Errorf("'%s' field %w", vErrs[0].StructField(), logger.ErrIsRequired)
	}
	g.validated = true
	return nil
}
