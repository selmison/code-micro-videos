package domain

import (
	"fmt"
	"mime/multipart"
	"strings"

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
	return [...]string{"Free", "10", "12", "14", "16", "18"}[*v]
}

func (v *VideoRating) Validate() error {
	switch *v {
	case FreeRating, TenRating, TwelveRating, FourteenRating, SixteenRating, EighteenRating:
		return nil
	}
	return fmt.Errorf("video rating %w", logger.ErrIsNotValidated)
}

type VideoValidatable struct {
	Id               string `validate:"not_blank"`
	Title            string `validate:"not_blank"`
	Description      string
	YearLaunched     *int16 `validate:"required"`
	Opened           bool
	Rating           *VideoRating          `validate:"required"`
	Duration         *int16                `validate:"required"`
	Categories       []CategoryValidatable `validate:"not_blank"`
	Genres           []GenreValidatable    `validate:"not_blank"`
	VideoFileHandler *multipart.FileHeader
	validated        bool
}

type Video struct {
	id               string
	title            string
	description      string
	yearLaunched     int16
	opened           bool
	rating           VideoRating
	duration         int16
	categories       []Category
	genres           []Genre
	videoFileHandler multipart.FileHeader
}

func NewVideo(vValidatable VideoValidatable) (*Video, error) {
	if err := vValidatable.Validate(); err != nil {
		return nil, err
	}
	vValidatable.Title = strings.TrimSpace(vValidatable.Title)
	vValidatable.Description = strings.TrimSpace(vValidatable.Description)
	categories, err := mapToCategories(vValidatable.Categories)
	if err != nil {
		return nil, err
	}
	genres, err := mapToGenres(vValidatable.Genres)
	if err != nil {
		return nil, err
	}
	video := &Video{
		id:               vValidatable.Id,
		title:            vValidatable.Title,
		description:      vValidatable.Description,
		yearLaunched:     *vValidatable.YearLaunched,
		opened:           vValidatable.Opened,
		rating:           *vValidatable.Rating,
		duration:         *vValidatable.Duration,
		categories:       categories,
		genres:           genres,
		videoFileHandler: *vValidatable.VideoFileHandler,
	}
	return video, nil
}

func mapToGenres(gValidatables []GenreValidatable) ([]Genre, error) {
	var genres []Genre
	if len(gValidatables) > 0 {
		genres = make([]Genre, len(gValidatables))
		for i, gValidatable := range gValidatables {
			genre, err := gValidatable.mapToGenre()
			if err != nil {
				return nil, err
			}
			genres[i] = *genre
		}
	}
	return genres, nil
}

func mapToCategories(cValidatables []CategoryValidatable) ([]Category, error) {
	var categories []Category
	if len(cValidatables) > 0 {
		categories = make([]Category, len(cValidatables))
		for i, cValidatable := range cValidatables {
			category, err := cValidatable.mapToCategory()
			if err != nil {
				return nil, err
			}
			categories[i] = *category
		}
	}
	return categories, nil
}

func (v *VideoValidatable) IsValidated() bool {
	return v.validated
}

func (v *VideoValidatable) Validate() error {
	v.validated = false
	err := validate.Struct(v)
	if err != nil {
		vErrs := err.(validator.ValidationErrors)
		return fmt.Errorf("'%s' field %w", vErrs[0].StructField(), logger.ErrIsRequired)
	}
	if err := v.Rating.Validate(); err != nil {
		return err
	}
	v.validated = true
	return nil
}

func (v *Video) MapToVideoValidatable() *VideoValidatable {
	id := v.id
	yearLaunched := v.yearLaunched
	rating := v.rating
	duration := v.duration
	videoFileHandler := v.videoFileHandler
	cValidatables := mapToCategoryValidatables(v.categories)
	gValidatables := mapToGenreValidatables(v.genres)
	return &VideoValidatable{
		Id:               id,
		Title:            v.title,
		Description:      v.description,
		YearLaunched:     &yearLaunched,
		Opened:           v.opened,
		Rating:           &rating,
		Duration:         &duration,
		Categories:       cValidatables,
		Genres:           gValidatables,
		VideoFileHandler: &videoFileHandler,
	}
}
