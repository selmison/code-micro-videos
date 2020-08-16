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
	"github.com/selmison/code-micro-videos/pkg/crud"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (r Repository) UpdateVideo(title string, videoDTO crud.VideoDTO) error {
	video, err := r.FetchVideo(title)
	if err != nil {
		return err
	}
	titleDTO := strings.ToLower(strings.TrimSpace(videoDTO.Title))
	video.Title = titleDTO
	_, err = video.UpdateG(r.ctx, boil.Infer())
	if err != nil {
		return fmt.Errorf("%s %w", titleDTO, logger.ErrAlreadyExists)
	}
	return nil
}

func (r Repository) AddVideo(videoDTO crud.VideoDTO) error {
	video := models.Video{
		ID:    uuid.New().String(),
		Title: strings.ToLower(strings.TrimSpace(videoDTO.Title)),
	}
	err := video.InsertG(r.ctx, boil.Infer())
	if err != nil {
		var e *pq.Error
		if errors.As(err, &e) {
			if e.Code.Name() == "unique_violation" {
				return fmt.Errorf("title '%s' %w", videoDTO.Title, logger.ErrAlreadyExists)
			} else {
				return fmt.Errorf("%s: %w", "method Repository.AddVideo(videoDTO)", err)
			}
		}
	}
	return nil
}

func (r Repository) RemoveVideo(title string) error {
	c, err := r.FetchVideo(title)
	if err != nil {
		return err
	}
	_, err = c.DeleteG(r.ctx, false)
	return err
}

func (r Repository) GetVideos(limit int) (models.VideoSlice, error) {
	if limit <= 0 {
		return nil, nil
	}
	videos, err := models.Videos(Limit(limit)).AllG(r.ctx)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (r Repository) FetchVideo(title string) (models.Video, error) {
	videoSlice, err := models.Videos(models.VideoWhere.Title.EQ(title)).AllG(r.ctx)
	if err != nil {
		return models.Video{}, err
	}
	if len(videoSlice) == 0 {
		return models.Video{}, sql.ErrNoRows
	}
	return *videoSlice[0], nil
}
