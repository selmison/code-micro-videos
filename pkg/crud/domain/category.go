package domain

import (
	"fmt"
	"strings"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

type Category struct {
	Id          string `validate:"not_blank"`
	Name        string `validate:"not_blank"`
	Description string
	Genres      []Genre
}

func (c *Category) Validate() error {
	if strings.ToLower(c.Name) != c.Name {
		return fmt.Errorf("lowercase 'Name' field %v", logger.ErrIsRequired)
	}
	err := validate.Struct(c)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		return fmt.Errorf("'%s' field %w", vErrs[0].StructField(), logger.ErrIsRequired)
	}
	return nil
}
