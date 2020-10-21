package domain

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

type Genre struct {
	Id         string `validate:"not_blank"`
	Name       string `validate:"not_blank"`
	Categories []Category
}

func (g *Genre) Validate() error {
	if strings.ToLower(g.Name) != g.Name {
		return fmt.Errorf("lowercase 'Name' field %v", logger.ErrIsRequired)
	}
	err := validate.Struct(g)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		return fmt.Errorf("'%s' field %w", vErrs[0].StructField(), logger.ErrIsRequired)
	}
	return nil
}
