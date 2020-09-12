package sqlboiler

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/pkg/storage/files"
)

type Repository struct {
	ctx       context.Context
	repoFiles files.Repository
}

func NewRepository(ctx context.Context, db *sql.DB, repoFiles files.Repository) *Repository {
	boil.SetDB(db)

	models.CreateCategoryHook(boil.BeforeInsertHook, isValidUUIDCategoryHook)
	models.CreateCategoryHook(boil.BeforeUpdateHook, isValidUUIDCategoryHook)
	models.CreateCategoryHook(boil.BeforeUpsertHook, isValidUUIDCategoryHook)

	models.CreateGenreHook(boil.BeforeInsertHook, isValidUUIDGenreHook)
	models.CreateGenreHook(boil.BeforeUpdateHook, isValidUUIDGenreHook)
	models.CreateGenreHook(boil.BeforeUpsertHook, isValidUUIDGenreHook)

	models.AddCastMemberHook(boil.BeforeInsertHook, isValidUUIDCastMemberHook)
	models.AddCastMemberHook(boil.BeforeUpdateHook, isValidUUIDCastMemberHook)
	models.AddCastMemberHook(boil.BeforeUpsertHook, isValidUUIDCastMemberHook)

	models.AddVideoHook(boil.BeforeInsertHook, isValidUUIDVideoHook)
	models.AddVideoHook(boil.BeforeUpdateHook, isValidUUIDVideoHook)
	models.AddVideoHook(boil.BeforeUpsertHook, isValidUUIDVideoHook)

	return &Repository{ctx, repoFiles}
}

func isValidUUIDCategoryHook(_ context.Context, _ boil.ContextExecutor, c *models.Category) error {
	if !isValidUUID(c.ID) {
		return fmt.Errorf("%s %w", "ID", logger.ErrIsNotValidated)
	}
	return nil
}

func isValidUUIDGenreHook(_ context.Context, _ boil.ContextExecutor, g *models.Genre) error {
	if !isValidUUID(g.ID) {
		return fmt.Errorf("%s %w", "ID", logger.ErrIsNotValidated)
	}
	return nil
}

func isValidUUIDCastMemberHook(_ context.Context, _ boil.ContextExecutor, c *models.CastMember) error {
	if !isValidUUID(c.ID) {
		return fmt.Errorf("%s %w", "ID", logger.ErrIsNotValidated)
	}
	return nil
}
func isValidUUIDVideoHook(_ context.Context, _ boil.ContextExecutor, v *models.Video) error {
	if !isValidUUID(v.ID) {
		return fmt.Errorf("%s %w", "ID", logger.ErrIsNotValidated)
	}
	return nil
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
