package sqlboiler

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

type Repository struct {
	ctx context.Context
}

func NewRepository(ctx context.Context, db *sql.DB) *Repository {
	boil.SetDB(db)
	models.AddCategoryHook(boil.BeforeInsertHook, isValidUUIDCategoryHook)
	models.AddCategoryHook(boil.BeforeUpdateHook, isValidUUIDCategoryHook)
	models.AddCategoryHook(boil.BeforeUpsertHook, isValidUUIDCategoryHook)

	models.AddGenreHook(boil.BeforeInsertHook, isValidUUIDGenreHook)
	models.AddGenreHook(boil.BeforeUpdateHook, isValidUUIDGenreHook)
	models.AddGenreHook(boil.BeforeUpsertHook, isValidUUIDGenreHook)

	models.AddCastMemberHook(boil.BeforeInsertHook, isValidUUIDCastMemberHook)
	models.AddCastMemberHook(boil.BeforeUpdateHook, isValidUUIDCastMemberHook)
	models.AddCastMemberHook(boil.BeforeUpsertHook, isValidUUIDCastMemberHook)

	return &Repository{ctx}
}
func isValidUUIDCategoryHook(ctx context.Context, exec boil.ContextExecutor, c *models.Category) error {
	if !isValidUUID(c.ID) {
		return fmt.Errorf("%s %w", "UUID", logger.ErrIsNotValidated)
	}
	return nil
}

func isValidUUIDGenreHook(ctx context.Context, exec boil.ContextExecutor, g *models.Genre) error {
	if !isValidUUID(g.ID) {
		return fmt.Errorf("%s %w", "UUID", logger.ErrIsNotValidated)
	}
	return nil
}

func isValidUUIDCastMemberHook(ctx context.Context, exec boil.ContextExecutor, c *models.CastMember) error {
	if !isValidUUID(c.ID) {
		return fmt.Errorf("%s %w", "UUID", logger.ErrIsNotValidated)
	}
	return nil
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
