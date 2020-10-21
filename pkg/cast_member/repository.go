package cast_member

//go:generate mockgen -destination=./mock/repository.go -package=mock . Repository

import (
	"context"
)

// Repository persists castMembers.
type Repository interface {
	// Store stores an castMember.
	Store(ctx context.Context, castMember CastMember) error

	// GetAll returns all castMembers.
	GetAll(ctx context.Context) ([]CastMember, error)

	// GetMany returns castMembers.
	GetMany(ctx context.Context, ids []string) ([]CastMember, error)

	// GetOne returns a single castMember by its Id.
	GetOne(ctx context.Context, id string) (CastMember, error)

	// DeleteOne deletes a single castMember by its Id.
	DeleteOne(ctx context.Context, id string) error

	// UpdateOne updates a single castMember by its Id.
	UpdateOne(ctx context.Context, id string, updateCastMember UpdateCastMemberDTO) error
}

//type repository struct {
//	repoFiles files.Repository
//}

//func NewRepository(db *sql.DB, repoFiles files.Repository) Repository {
//boil.SetDB(db)
//return &repository{repoFiles}
//}
