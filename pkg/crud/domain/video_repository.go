package domain

import "github.com/google/uuid"

type VideoRepository interface {
	CreateVideo(ctx Context, video Video) (uuid.UUID, error)
	FetchVideo(ctx Context, name string) (Video, error)
	GetVideos(ctx Context, limit int) ([]Video, error)
	RemoveVideo(ctx Context, name string) error
	UpdateVideo(ctx Context, name string, video Video) error
}
