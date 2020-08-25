package crud

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

var categoryValidate *validator.Validate

type CategoryDTO struct {
	Name        string     `json:"name" validate:"not_blank"`
	Description string     `json:"description,omitempty"`
	Genres      []GenreDTO `json:"genres" validate:"not_blank"`
}

func MapGenreToDTO(category models.Category) (*CategoryDTO, error) {
	genreDTOs := make([]GenreDTO, len(category.R.Genres))
	for i, genre := range category.R.Genres {
		genreDTOs[i] = GenreDTO{
			Name: genre.Name,
		}
	}
	dto := &CategoryDTO{
		Name:   category.Name,
		Genres: genreDTOs,
	}
	if err := dto.Validate(); err != nil {
		return nil, err
	}
	return dto, nil
}

func (c *CategoryDTO) Validate() error {
	err := categoryValidate.Struct(c)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		return fmt.Errorf("'%s' field %w", vErrs[0].StructField(), logger.ErrIsRequired)
	}
	return nil
}

func init() {
	categoryValidate = validator.New()
	categoryValidate.RegisterValidation("not_blank", validators.NotBlank)
}
