package video

//go:generate mockgen -destination=./mock/repository.go -package=mock . Repository

import "context"

// Repository persists videos.
type Repository interface {
	// Store stores an video.
	Store(ctx context.Context, video Video) error

	// GetAll returns all videos.
	GetAll(ctx context.Context) ([]Video, error)

	// GetMany returns videos.
	GetMany(ctx context.Context, ids []string) ([]Video, error)

	// GetOne returns a single video by its Id.
	GetOne(ctx context.Context, id string) (Video, error)

	// DeleteOne deletes a single video by its Id.
	DeleteOne(ctx context.Context, id string) error

	// UpdateOne updates a single video by its Id.
	UpdateOne(ctx context.Context, id string, updateVideo UpdateVideo) error
}
