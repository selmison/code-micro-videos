//go:generate mockgen -destination=./mock/service.go -package=mock . Service

package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/selmison/code-micro-videos/pkg/crud/domain"
)

type service struct {
	r domain.Repository
}

type Service interface {
	GetCategories(ctx context.Context, limit int) ([]domain.Category, error)
	FetchCategory(ctx context.Context, name string) (domain.Category, error)
	CreateCategory(ctx context.Context, fields domain.GenreValidatable) error
	RemoveCategory(ctx context.Context, name string) error
	UpdateCategory(ctx context.Context, name string, fields domain.GenreValidatable) error

	GetCastMembers(ctx context.Context, limit int) ([]domain.CastMember, error)
	FetchCastMember(ctx context.Context, name string) (domain.CastMember, error)
	AddCastMember(ctx context.Context, fields domain.GenreValidatable) error
	RemoveCastMember(ctx context.Context, name string) error
	UpdateCastMember(ctx context.Context, name string, fields domain.GenreValidatable) error

	GetGenres(ctx context.Context, limit int) ([]domain.Genre, error)
	FetchGenre(ctx context.Context, name string) (domain.Genre, error)
	CreateGenre(ctx context.Context, fields domain.GenreValidatable) error
	RemoveGenre(ctx context.Context, name string) error
	UpdateGenre(ctx context.Context, name string, fields domain.GenreValidatable) error

	GetVideos(ctx context.Context, limit int) ([]domain.Video, error)
	FetchVideo(ctx context.Context, title string) (domain.Video, error)
	CreateVideo(ctx context.Context, fields domain.VideoValidatable) (uuid.UUID, error)
	RemoveVideo(ctx context.Context, title string) error
	UpdateVideo(ctx context.Context, title string, fields domain.VideoValidatable) (uuid.UUID, error)
}

// NewService creates a crud service with the necessary dependencies
func NewService(repoDB domain.Repository) *service {
	return &service{repoDB}
}
