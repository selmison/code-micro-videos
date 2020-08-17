package crud

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (s service) RemoveVideo(name string) error {
	if len(strings.TrimSpace(name)) == 0 {
		return fmt.Errorf("'name' %w", logger.ErrIsRequired)
	}
	if err := s.r.RemoveVideo(name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", name, logger.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s service) UpdateVideo(title string, videoDTO VideoDTO) error {
	if len(strings.TrimSpace(title)) == 0 {
		return fmt.Errorf("'title' %w", logger.ErrIsRequired)
	}
	if err := videoDTO.Validate(); err != nil {
		return err
	}
	if err := s.r.UpdateVideo(title, videoDTO); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", title, logger.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s service) AddVideo(videoDTO VideoDTO) error {
	if err := videoDTO.Validate(); err != nil {
		return err
	}
	return s.r.AddVideo(videoDTO)
}

func (s service) GetVideos(limit int) (models.VideoSlice, error) {
	if limit < 0 {
		return nil, logger.ErrInvalidedLimit
	}
	return s.r.GetVideos(limit)
}

func (s service) FetchVideo(name string) (models.Video, error) {
	name = strings.TrimSpace(name)
	c, err := s.r.FetchVideo(name)
	if err == sql.ErrNoRows {
		return models.Video{}, fmt.Errorf("%s: %w", name, logger.ErrNotFound)
	} else if err != nil {
		return models.Video{}, fmt.Errorf("%s: %w", name, logger.ErrInternalApplication)
	}

	return c, nil
}
