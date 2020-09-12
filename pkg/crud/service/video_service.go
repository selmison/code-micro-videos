package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/selmison/code-micro-videos/pkg/crud/domain"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (s service) CreateVideo(ctx domain.Context, fields domain.VideoValidatable) (uuid.UUID, error) {
	video, err := domain.NewVideo(fields)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error CreateVideo(): %w", err)
	}
	return s.r.CreateVideo(ctx, *video)
}

func (s service) FetchVideo(ctx domain.Context, title string) (domain.Video, error) {
	title = strings.ToLower(strings.TrimSpace(title))
	if len(title) == 0 {
		return domain.Video{}, fmt.Errorf("'title' %w", logger.ErrIsRequired)
	}
	video, err := s.r.FetchVideo(ctx, title)
	if err == sql.ErrNoRows {
		return domain.Video{}, fmt.Errorf("%s: %w", title, logger.ErrNotFound)
	} else if err != nil {
		return domain.Video{}, fmt.Errorf("%s: %w", title, logger.ErrInternalApplication)
	}
	return video, nil
}

func (s service) GetVideos(ctx domain.Context, limit int) ([]domain.Video, error) {
	if limit < 0 {
		return nil, logger.ErrInvalidedLimit
	}
	return s.r.GetVideos(ctx, limit)
}

func (s service) RemoveVideo(ctx domain.Context, title string) error {
	title = strings.ToLower(strings.TrimSpace(title))
	if len(title) == 0 {
		return fmt.Errorf("'title' %w", logger.ErrIsRequired)
	}
	if err := s.r.RemoveVideo(ctx, title); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", title, logger.ErrNotFound)
		}
		return err
	}
	return nil
}

func (s service) UpdateVideo(ctx domain.Context, title string, fields domain.VideoValidatable) error {
	title = strings.ToLower(strings.TrimSpace(title))
	video, err := domain.NewVideo(fields)
	if err != nil {
		return fmt.Errorf("error UpdateVideo(): %w", err)
	}
	if err := s.r.UpdateVideo(ctx, title, *video); err != nil {
		return err
	}
	return nil
}
