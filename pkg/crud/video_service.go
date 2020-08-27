package crud

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

func (s service) UpdateVideo(title string, videoDTO VideoDTO) error {
	title = strings.ToLower(strings.TrimSpace(title))
	if len(title) == 0 {
		return fmt.Errorf("'title' %w", logger.ErrIsRequired)
	}
	if err := videoDTO.Validate(); err != nil {
		return err
	}
	videoDTO.Title = strings.ToLower(strings.TrimSpace(videoDTO.Title))
	videoDTO.Description = strings.TrimSpace(videoDTO.Description)
	if err := s.r.UpdateVideo(title, videoDTO); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", title, logger.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s service) AddVideo(videoDTO VideoDTO) error {
	videoDTO.Title = strings.ToLower(strings.TrimSpace(videoDTO.Title))
	videoDTO.Description = strings.TrimSpace(videoDTO.Description)
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
