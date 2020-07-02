package modifying

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

var validate *validator.Validate

type CategoryDTO struct {
	Name        string `json:"name" validate:"not_blank"`
	Description string `json:"description,omitempty"`
}

func (c *CategoryDTO) Validate() error {
	err := validate.Struct(c)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		return fmt.Errorf("%s %w", vErrs[0].StructField(), ErrIsRequired)
	}
	return nil
}

func init() {
	validate = validator.New()
	validate.RegisterValidation("not_blank", validators.NotBlank)
}
