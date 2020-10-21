package genre

//go:generate mockgen -destination=./mock/repository.go -package=mock . Repository

import "context"

// Repository persists genres.
type Repository interface {
	// Store stores an genre.
	Store(ctx context.Context, genre Genre) error

	// GetAll returns all genres.
	GetAll(ctx context.Context) ([]Genre, error)

	// GetMany returns genres.
	GetMany(ctx context.Context, ids []string) ([]Genre, error)

	// GetOne returns a single genre by its Id.
	GetOne(ctx context.Context, id string) (Genre, error)

	// DeleteOne deletes a single genre by its Id.
	DeleteOne(ctx context.Context, id string) error

	// UpdateOne updates a single genre by its Id.
	UpdateOne(ctx context.Context, id string, updateGenre UpdateGenre) error
}
