package crud

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

var videoValidate *validator.Validate

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
	return [...]string{"Free", "10", "12", "14", "16", "18"}[*v]
}

func (v *VideoRating) Validate() error {
	switch *v {
	case FreeRating, TenRating, TwelveRating, FourteenRating, SixteenRating, EighteenRating:
		return nil
	}
	return fmt.Errorf("video rating %w", logger.ErrIsNotValidated)
}

type VideoDTO struct {
	Title        string        `json:"title" validate:"not_blank"`
	Description  string        `json:"description"`
	YearLaunched *int16        `json:"year_launched" validate:"required"`
	Opened       bool          `json:"opened"`
	Rating       *VideoRating  `json:"rating" validate:"required"`
	Duration     *int16        `json:"duration" validate:"required"`
	VideoFile    []byte        `json:"video_file"`
	Categories   []CategoryDTO `json:"categories" validate:"not_blank"`
	Genres       []GenreDTO    `json:"genres" validate:"not_blank"`
}

func MapVideoToDTO(video models.Video) (*VideoDTO, error) {
	categoriesDTOs := make([]CategoryDTO, len(video.R.Categories))
	for i, category := range video.R.Categories {
		categoriesDTOs[i] = CategoryDTO{
			Name:        category.Name,
			Description: category.Description.String,
		}
	}
	genresDTOs := make([]GenreDTO, len(video.R.Genres))
	for i, genre := range video.R.Genres {
		genresDTOs[i] = GenreDTO{
			Name: genre.Name,
		}
	}
	rating := VideoRating(video.Rating)
	dto := &VideoDTO{
		Title:        video.Title,
		Description:  video.Description,
		YearLaunched: &video.YearLaunched,
		Opened:       video.Opened.Bool,
		Rating:       &rating,
		Duration:     &video.Duration,
		Categories:   categoriesDTOs,
		Genres:       genresDTOs,
	}
	if err := dto.Validate(); err != nil {
		return nil, err
	}
	return dto, nil
}

func (v *VideoDTO) Validate() error {
	err := videoValidate.Struct(v)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		return fmt.Errorf("'%s' field %w", vErrs[0].StructField(), logger.ErrIsRequired)
	}
	if err := v.Rating.Validate(); err != nil {
		return err
	}
	return nil
}

func init() {
	videoValidate = validator.New()
	videoValidate.RegisterValidation("not_blank", validators.NotBlank)
}
