package sqlboiler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/cast_member"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

// Store keeps castMembers in the SqlBoiler.
type Store struct{}

func (s *Store) UpdateOne(ctx context.Context, id string, updateCastMember cast_member.CastMember) error {
	castMember, err := getById(ctx, id)
	if err != nil {
		return err
	}
	nameDTO := strings.ToLower(strings.TrimSpace(updateCastMember.Name()))
	castMember.Name = nameDTO
	_, err = castMember.UpdateG(ctx, boil.Infer())
	if err != nil {
		return fmt.Errorf("%s %w", nameDTO, logger.ErrAlreadyExists)
	}
	return nil
}

func getById(ctx context.Context, id string) (models.CastMember, error) {
	castMemberSlice, err := models.CastMembers(models.CastMemberWhere.ID.EQ(id)).AllG(ctx)
	if err != nil {
		return models.CastMember{}, err
	}
	if len(castMemberSlice) == 0 {
		return models.CastMember{}, sql.ErrNoRows
	}
	return *castMemberSlice[0], nil
}

func (s *Store) Store(ctx context.Context, newCastMember cast_member.CastMember) error {
	castMember := models.CastMember{
		ID:   newCastMember.Id(),
		Name: newCastMember.Name(),
	}
	err := castMember.InsertG(ctx, boil.Infer())
	if err != nil {
		var e *pq.Error
		if errors.As(err, &e) {
			if e.Code.Name() == "unique_violation" {
				return fmt.Errorf("name '%s' %w", castMember.Name, logger.ErrAlreadyExists)
			} else {
				return fmt.Errorf("%s: %w", "method Repository.AddCastMember(castMemberDTO)", err)
			}
		}
	}
	return nil
}

func (s *Store) DeleteOne(ctx context.Context, id string) error {
	c, err := getById(ctx, id)
	if err != nil {
		return err
	}
	_, err = c.DeleteG(ctx, false)
	return err
}

func (s *Store) GetAll(ctx context.Context) ([]cast_member.CastMember, error) {
	castMemberSqlBoilers, err := models.CastMembers().AllG(ctx)
	if err != nil {
		return nil, err
	}
	castMembers := make([]cast_member.CastMember, len(castMemberSqlBoilers))
	for i, castMemberSqlBoiler := range castMemberSqlBoilers {
		castMembers[i], err = cast_member.NewCastMember(
			castMemberSqlBoiler.ID,
			cast_member.NewCastMemberDTO{
				Name: castMemberSqlBoiler.Name,
				Type: cast_member.CastMemberType(castMemberSqlBoiler.Type),
			},
		)
		if err != nil {
			return nil, err
		}
	}
	return castMembers, nil
}

func (s *Store) GetOne(ctx context.Context, id string) (cast_member.CastMember, error) {
	castMemberSlice, err := models.CastMembers(models.CastMemberWhere.Name.EQ(id)).AllG(ctx)
	if err != nil {
		return nil, err
	}
	if len(castMemberSlice) == 0 {
		return nil, sql.ErrNoRows
	}
	return cast_member.NewCastMember(
		castMemberSlice[0].ID,
		cast_member.NewCastMemberDTO{
			Name: castMemberSlice[0].Name,
			Type: cast_member.CastMemberType(castMemberSlice[0].Type),
		},
	)
}
