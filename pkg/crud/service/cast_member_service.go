package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/selmison/code-micro-videos/pkg/crud/domain"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (s service) CreateCastMember(ctx domain.Context, castMember domain.CastMember) error {
	if err := castMember.Validate(); err != nil {
		return err
	}
	return s.r.CreateCastMember(ctx, castMember)
}

func (s service) FetchCastMember(ctx domain.Context, name string) (domain.CastMember, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	if len(name) == 0 {
		return domain.CastMember{}, fmt.Errorf("'name' %w", logger.ErrIsRequired)
	}
	castMember, err := s.r.FetchCastMember(ctx, name)
	if err == sql.ErrNoRows {
		return domain.CastMember{}, fmt.Errorf("%s: %w", name, logger.ErrNotFound)
	} else if err != nil {
		return domain.CastMember{}, fmt.Errorf("%s: %w", name, logger.ErrInternalApplication)
	}
	return castMember, nil
}

func (s service) GetCastMembers(ctx domain.Context, limit int) ([]domain.CastMember, error) {
	if limit < 0 {
		return nil, logger.ErrInvalidedLimit
	}
	return s.r.GetCastMembers(ctx, limit)
}

func (s service) RemoveCastMember(ctx domain.Context, name string) error {
	name = strings.ToLower(strings.TrimSpace(name))
	if len(name) == 0 {
		return fmt.Errorf("'name' %w", logger.ErrIsRequired)
	}
	if err := s.r.RemoveCastMember(ctx, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", name, logger.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s service) UpdateCastMember(ctx domain.Context, name string, castMember domain.CastMember) error {
	name = strings.ToLower(strings.TrimSpace(name))
	if err := castMember.Validate(); err != nil {
		return err
	}
	if err := s.r.UpdateCastMember(ctx, name, castMember); err != nil {
		return err
	}
	return nil
}
