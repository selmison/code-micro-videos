package crud

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

var genreValidate *validator.Validate

type GenreDTO struct {
	Name       string        `json:"name" validate:"not_blank"`
	Categories []CategoryDTO `json:"categories" validate:"not_blank"`
}

func MapGenreToDTO(genre models.Genre) (*GenreDTO, error) {
	categoryDTOs := make([]CategoryDTO, len(genre.R.Categories))
	for i, genre := range genre.R.Categories {
		categoryDTOs[i] = CategoryDTO{
			Name: genre.Name,
		}
	}
	dto := &GenreDTO{
		Name:       genre.Name,
		Categories: categoryDTOs,
	}
	if err := dto.Validate(); err != nil {
		return nil, err
	}
	return dto, nil
}

func (c *GenreDTO) Validate() error {
	err := genreValidate.Struct(c)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		return fmt.Errorf("'%s' field %w", vErrs[0].StructField(), logger.ErrIsRequired)
	}
	return nil
}

func init() {
	genreValidate = validator.New()
	genreValidate.RegisterValidation("not_blank", validators.NotBlank)
}
