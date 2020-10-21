package category

//go:generate mockgen -destination=./mock/repository.go -package=mock . Repository

import "context"

// Repository persists categories.
type Repository interface {
	// Store stores an category.
	Store(ctx context.Context, category Category) error

	// GetAll returns all categories.
	GetAll(ctx context.Context) ([]Category, error)

	// GetMany returns categories.
	GetMany(ctx context.Context, ids []string) ([]Category, error)

	// GetOne returns a single category by its Id.
	GetOne(ctx context.Context, id string) (Category, error)

	// DeleteOne deletes a single category by its Id.
	DeleteOne(ctx context.Context, id string) error

	// UpdateOne updates a single category by its Id.
	UpdateOne(ctx context.Context, id string, updateCategory UpdateCategory) error
}
