package crud

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

var castMemberValidate *validator.Validate

type CastMemberType int16

const (
	Director CastMemberType = iota
	Actor
)

func (c CastMemberType) String() string {
	return [...]string{"Director", "Actor"}[c]
}

type CastMemberDTO struct {
	Name string         `json:"name" validate:"not_blank"`
	Type CastMemberType `json:"type"`
}

func (c CastMemberType) Validate() error {
	switch c {
	case Director, Actor:
		return nil
	}
	return fmt.Errorf("cast member type %w", logger.ErrIsNotValidated)
}

func (c *CastMemberDTO) Validate() error {
	err := castMemberValidate.Struct(c)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		return fmt.Errorf("%s %w", vErrs[0].StructField(), logger.ErrIsRequired)
	}
	if err := c.Type.Validate(); err != nil {
		return err
	}
	return nil
}

func init() {
	castMemberValidate = validator.New()
	castMemberValidate.RegisterValidation("not_blank", validators.NotBlank)
}
