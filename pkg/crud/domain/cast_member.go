package domain

import (
	"fmt"
	"log"
	"strings"

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

func (c CastMemberType) Validate() error {
	switch c {
	case Director, Actor:
		return nil
	}
	return fmt.Errorf("cast member type %w", logger.ErrIsNotValidated)
}

type CastMemberValidatable struct {
	Id        string `validate:"not_blank"`
	Name      string `validate:"not_blank"`
	Type      CastMemberType
	validated bool
}

type CastMember struct {
	id             string
	name           string
	castMemberType CastMemberType
}

func NewCastMember(fields CastMemberValidatable) (*CastMember, error) {
	if err := fields.Validate(); err != nil {
		return nil, err
	}
	fields.Name = strings.TrimSpace(fields.Name)
	castMember := &CastMember{
		id:             fields.Id,
		name:           fields.Name,
		castMemberType: fields.Type,
	}
	return castMember, nil
}

func (c *CastMember) MapToCastMemberValidatable() *CastMemberValidatable {
	return &CastMemberValidatable{
		Id:        c.id,
		Name:      c.name,
		Type:      c.castMemberType,
		validated: true,
	}
}

func (c *CastMemberValidatable) IsValidated() bool {
	return c.validated
}

func (c *CastMemberValidatable) Validate() error {
	c.validated = false
	err := castMemberValidate.Struct(c)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		return fmt.Errorf("%s %w", vErrs[0].StructField(), logger.ErrIsRequired)
	}
	if err := c.Type.Validate(); err != nil {
		return err
	}
	c.validated = true
	return nil
}

func init() {
	castMemberValidate = validator.New()
	if err := castMemberValidate.RegisterValidation("not_blank", validators.NotBlank); err != nil {
		log.Fatalln(err)
	}
}
