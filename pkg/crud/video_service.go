package crud

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (s service) RemoveVideo(title string) error {
	title = strings.ToLower(strings.TrimSpace(title))
	if len(title) == 0 {
		return fmt.Errorf("'title' %w", logger.ErrIsRequired)
	}
	if err := s.r.RemoveVideo(title); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", title, logger.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s service) UpdateVideo(title string, videoDTO VideoDTO) (uuid.UUID, error) {
	title = strings.ToLower(strings.TrimSpace(title))
	if len(title) == 0 {
		return uuid.UUID{}, fmt.Errorf("'title' %w", logger.ErrIsRequired)
	}
	if err := videoDTO.Validate(); err != nil {
		return uuid.UUID{}, err
	}
	videoDTO.Title = strings.ToLower(strings.TrimSpace(videoDTO.Title))
	videoDTO.Description = strings.TrimSpace(videoDTO.Description)
	id, err := s.r.UpdateVideo(title, videoDTO)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, fmt.Errorf("%s: %w", title, logger.ErrNotFound)
		}
		return uuid.UUID{}, err
	}
	return id, nil
}

func (s service) AddVideo(videoDTO VideoDTO) (uuid.UUID, error) {
	videoDTO.Title = strings.ToLower(strings.TrimSpace(videoDTO.Title))
	videoDTO.Description = strings.TrimSpace(videoDTO.Description)
	if err := videoDTO.Validate(); err != nil {
		return uuid.UUID{}, err
	}
	id, err := s.r.AddVideo(videoDTO)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (s service) GetVideos(limit int) (models.VideoSlice, error) {
	if limit < 0 {
		return nil, logger.ErrInvalidedLimit
	}
	return s.r.GetVideos(limit)
}

func (s service) FetchVideo(title string) (models.Video, error) {
	title = strings.ToLower(strings.TrimSpace(title))
	c, err := s.r.FetchVideo(title)
	if err == sql.ErrNoRows {
		return models.Video{}, fmt.Errorf("%s: %w", title, logger.ErrNotFound)
	} else if err != nil {
		return models.Video{}, fmt.Errorf("%s: %w", title, logger.ErrInternalApplication)
	}

	return c, nil
}
