package sqlboiler

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud/service"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (r Repository) UpdateCastMember(name string, castMemberDTO service.CastMemberDTO) error {
	castMember, err := r.FetchCastMember(name)
	if err != nil {
		return err
	}
	nameDTO := strings.ToLower(strings.TrimSpace(castMemberDTO.Name))
	castMember.Name = nameDTO
	_, err = castMember.UpdateG(r.ctx, boil.Infer())
	if err != nil {
		return fmt.Errorf("%s %w", nameDTO, logger.ErrAlreadyExists)
	}
	return nil
}

func (r Repository) AddCastMember(castMemberDTO service.CastMemberDTO) error {
	castMember := models.CastMember{
		ID:   uuid.New().String(),
		Name: strings.ToLower(strings.TrimSpace(castMemberDTO.Name)),
	}
	err := castMember.InsertG(r.ctx, boil.Infer())
	if err != nil {
		var e *pq.Error
		if errors.As(err, &e) {
			if e.Code.Name() == "unique_violation" {
				return fmt.Errorf("name '%s' %w", castMemberDTO.Name, logger.ErrAlreadyExists)
			} else {
				return fmt.Errorf("%s: %w", "method Repository.AddCastMember(castMemberDTO)", err)
			}
		}
	}
	return nil
}

func (r Repository) RemoveCastMember(name string) error {
	c, err := r.FetchCastMember(name)
	if err != nil {
		return err
	}
	_, err = c.DeleteG(r.ctx, false)
	return err
}

func (r Repository) GetCastMembers(limit int) (models.CastMemberSlice, error) {
	if limit <= 0 {
		return nil, nil
	}
	castMembers, err := models.CastMembers(Limit(limit)).AllG(r.ctx)
	if err != nil {
		return nil, err
	}
	return castMembers, nil
}

func (r Repository) FetchCastMember(name string) (models.CastMember, error) {
	castMemberSlice, err := models.CastMembers(models.CastMemberWhere.Name.EQ(name)).AllG(r.ctx)
	if err != nil {
		return models.CastMember{}, err
	}
	if len(castMemberSlice) == 0 {
		return models.CastMember{}, sql.ErrNoRows
	}
	return *castMemberSlice[0], nil
}
