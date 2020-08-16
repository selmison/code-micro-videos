package crud

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

var videoValidate *validator.Validate

type VideoRating int16

const (
	FreeRating VideoRating = iota
	TenRating
	TwelveRating
	FourteenRating
	SixteenRating
	EighteenRating
)

func (v VideoRating) String() string {
	return [...]string{"Free", "10", "12", "14", "16", "18"}[v]
}

func (v VideoRating) Validate() error {
	switch v {
	case FreeRating, TenRating, TwelveRating, FourteenRating, SixteenRating, EighteenRating:
		return nil
	}
	return fmt.Errorf("video rating %w", logger.ErrIsNotValidated)
}

type VideoDTO struct {
	Title        string      `json:"title" validate:"not_blank"`
	Description  string      `json:"description"`
	YearLaunched int16       `json:"year_launched"`
	Opened       bool        `json:"opened"`
	Rating       VideoRating `json:"rating"`
	Duration     int16       `json:"duration"`
}

func (c *VideoDTO) Validate() error {
	err := videoValidate.Struct(c)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		return fmt.Errorf("%s %w", vErrs[0].StructField(), logger.ErrIsRequired)
	}
	if err := c.Rating.Validate(); err != nil {
		return err
	}
	return nil
}

func init() {
	videoValidate = validator.New()
	videoValidate.RegisterValidation("not_blank", validators.NotBlank)
}
