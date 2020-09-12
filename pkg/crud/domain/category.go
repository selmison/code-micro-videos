package domain

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

type CategoryValidatable struct {
	Id          string `validate:"not_blank"`
	Name        string `validate:"required"`
	Description string
	Genres      []GenreValidatable
	validated   bool
}

type Category struct {
	id          string
	name        string
	description string
	genres      []Genre
}

func (c *Category) Name() string {
	return c.name
}

func NewCategory(fields CategoryValidatable) (*Category, error) {
	return fields.mapToCategory()
}

func (c *Category) MapToCategoryValidatable() *CategoryValidatable {
	var genres []GenreValidatable
	if c.genres != nil && len(c.genres) > 0 {
		genres = make([]GenreValidatable, len(c.genres))
		for i, category := range c.genres {
			genres[i] = *category.MapToGenreValidatable()
		}
	}
	return &CategoryValidatable{
		Id:          c.id,
		Name:        c.name,
		Description: c.description,
		Genres:      genres,
		validated:   true,
	}
}

func mapToCategoryValidatables(categories []Category) []CategoryValidatable {
	var cValidatables []CategoryValidatable
	if len(categories) > 0 {
		cValidatables = make([]CategoryValidatable, len(categories))
		for i, category := range categories {
			cValidatables[i] = *category.MapToCategoryValidatable()
		}
	}
	return cValidatables
}

func (c *CategoryValidatable) IsValidated() bool {
	return c.validated
}

func (c *CategoryValidatable) mapToCategory() (*Category, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	var genres []Genre
	if len(c.Genres) > 0 {
		genres = make([]Genre, len(c.Genres))
		for i, categoryValidatable := range c.Genres {
			genre, err := categoryValidatable.mapToGenre()
			if err != nil {
				return nil, err
			}
			genres[i] = *genre
		}
	}
	return &Category{
		id:          c.Id,
		name:        c.Name,
		description: c.Description,
		genres:      genres,
	}, nil
}

func (c *CategoryValidatable) Validate() error {
	c.validated = false
	if strings.ToLower(c.Name) != c.Name {
		return fmt.Errorf("lowercase 'Name' field %v\n", logger.ErrIsRequired)
	}
	err := validate.Struct(c)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		return fmt.Errorf("'%s' field %w", vErrs[0].StructField(), logger.ErrIsRequired)
	}
	c.validated = true
	return nil
}
