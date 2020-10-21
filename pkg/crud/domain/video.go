package domain

import (
	"fmt"
	"mime/multipart"

	"github.com/go-playground/validator/v10"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

type VideoRating int16

const (
	FreeRating VideoRating = iota + 1
	TenRating
	TwelveRating
	FourteenRating
	SixteenRating
	EighteenRating
)

func (v *VideoRating) String() string {
	return [...]string{"Free", "10", "12", "14", "16", "18"}[*v-1]
}

func (v *VideoRating) Validate() error {
	switch *v {
	case FreeRating, TenRating, TwelveRating, FourteenRating, SixteenRating, EighteenRating:
		return nil
	}
	return fmt.Errorf("video rating %w", logger.ErrIsNotValidated)
}

type Video struct {
	Id               string `validate:"not_blank"`
	Title            string `validate:"not_blank"`
	Description      string
	YearLaunched     *int16 `validate:"required"`
	Opened           bool
	Rating           *VideoRating `validate:"required"`
	Duration         *int16       `validate:"required"`
	Categories       []Category   `validate:"not_blank"`
	Genres           []Genre      `validate:"not_blank"`
	VideoFileHandler *multipart.FileHeader
}

func (v *Video) Validate() error {
	err := validate.Struct(v)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		return fmt.Errorf("'%s' field %w", vErrs[0].StructField(), logger.ErrIsRequired)
	}
	if err := v.Rating.Validate(); err != nil {
		return err
	}
	return nil
}
