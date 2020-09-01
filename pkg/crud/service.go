//go:generate mockgen -destination=./mock/service.go -package=mock . Repository,Service

package crud

import (
	"github.com/google/uuid"

	"github.com/selmison/code-micro-videos/models"
)

type Repository interface {
	Service
}

type service struct {
	r Repository
}

type Service interface {
	GetCategories(limit int) (models.CategorySlice, error)
	FetchCategory(name string) (models.Category, error)
	AddCategory(dto CategoryDTO) error
	RemoveCategory(name string) error
	UpdateCategory(name string, dto CategoryDTO) error

	GetCastMembers(limit int) (models.CastMemberSlice, error)
	FetchCastMember(name string) (models.CastMember, error)
	AddCastMember(dto CastMemberDTO) error
	RemoveCastMember(name string) error
	UpdateCastMember(name string, dto CastMemberDTO) error

	GetGenres(limit int) (models.GenreSlice, error)
	FetchGenre(name string) (models.Genre, error)
	AddGenre(dto GenreDTO) error
	RemoveGenre(name string) error
	UpdateGenre(name string, dto GenreDTO) error

	GetVideos(limit int) (models.VideoSlice, error)
	FetchVideo(name string) (models.Video, error)
	AddVideo(dto VideoDTO) (uuid.UUID, error)
	RemoveVideo(name string) error
	UpdateVideo(name string, dto VideoDTO) (uuid.UUID, error)
}

// NewService creates a crud service with the necessary dependencies
func NewService(repoDB Repository) *service {
	return &service{repoDB}
}
