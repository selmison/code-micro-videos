package crud

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

var categoryValidate *validator.Validate

type CategoryDTO struct {
	Name        string `json:"name" validate:"not_blank"`
	Description string `json:"description,omitempty"`
}

func (c *CategoryDTO) Validate() error {
	err := categoryValidate.Struct(c)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		return fmt.Errorf("%s %w", vErrs[0].StructField(), logger.ErrIsRequired)
	}
	return nil
}

func init() {
	categoryValidate = validator.New()
	categoryValidate.RegisterValidation("not_blank", validators.NotBlank)
}
