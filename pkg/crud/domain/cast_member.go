package domain

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

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

type CastMember struct {
	Id   string `validate:"not_blank"`
	Name string `validate:"not_blank"`
	Type CastMemberType
}

func (c *CastMember) Validate() error {
	err := validate.Struct(c)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		return fmt.Errorf("'%s' field %w", vErrs[0].StructField(), logger.ErrIsRequired)
	}
	if err := c.Type.Validate(); err != nil {
		return err
	}
	return nil
}
