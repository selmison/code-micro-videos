package crud

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

var genreValidate *validator.Validate

type GenreDTO struct {
	Name string `json:"name" validate:"not_blank"`
}

func (c *GenreDTO) Validate() error {
	err := genreValidate.Struct(c)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		return fmt.Errorf("%s %w", vErrs[0].StructField(), logger.ErrIsRequired)
	}
	return nil
}

func init() {
	genreValidate = validator.New()
	genreValidate.RegisterValidation("not_blank", validators.NotBlank)
}
