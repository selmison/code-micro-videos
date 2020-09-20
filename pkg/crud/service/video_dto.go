package service

import (
	"mime/multipart"

	"github.com/selmison/code-micro-videos/pkg/crud/domain"
)

type VideoDTO struct {
	Title        string                `json:"title" schema:"title" validate:"not_blank"`
	Description  string                `json:"description" schema:"description"`
	YearLaunched *int16                `json:"year_launched" schema:"year_launched" validate:"required"`
	Opened       bool                  `json:"opened" schema:"opened"`
	Rating       *domain.VideoRating   `json:"rating" schema:"rating" validate:"required"`
	Duration     *int16                `json:"duration" schema:"duration" validate:"required"`
	Categories   []domain.Category     `json:"categories" schema:"categories" validate:"not_blank"`
	Genres       []domain.Genre        `json:"genres" schema:"genres" validate:"not_blank"`
	File         *multipart.FileHeader `json:"-" schema:"-"`
}
