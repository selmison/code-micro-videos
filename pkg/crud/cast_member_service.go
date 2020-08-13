package crud

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (s service) RemoveCastMember(name string) error {
	if len(strings.TrimSpace(name)) == 0 {
		return fmt.Errorf("'name' %w", logger.ErrIsRequired)
	}
	if err := s.r.RemoveCastMember(name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", name, logger.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s service) UpdateCastMember(name string, castMemberDTO CastMemberDTO) error {
	if len(strings.TrimSpace(name)) == 0 {
		return fmt.Errorf("'name' %w", logger.ErrIsRequired)
	}
	if err := castMemberDTO.Validate(); err != nil {
		return err
	}
	if err := s.r.UpdateCastMember(name, castMemberDTO); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", name, logger.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s service) AddCastMember(castMemberDTO CastMemberDTO) error {
	castMemberDTO.Name = strings.TrimSpace(castMemberDTO.Name)
	if err := castMemberDTO.Validate(); err != nil {
		return err
	}
	return s.r.AddCastMember(castMemberDTO)
}

func (s service) GetCastMembers(limit int) (models.CastMemberSlice, error) {
	if limit < 0 {
		return nil, logger.ErrInvalidedLimit
	}
	return s.r.GetCastMembers(limit)
}

func (s service) FetchCastMember(name string) (models.CastMember, error) {
	name = strings.TrimSpace(name)
	c, err := s.r.FetchCastMember(name)
	if err == sql.ErrNoRows {
		return models.CastMember{}, fmt.Errorf("%s: %w", name, logger.ErrNotFound)
	} else if err != nil {
		return models.CastMember{}, fmt.Errorf("%s: %w", name, logger.ErrInternalApplication)
	}

	return c, nil
}
