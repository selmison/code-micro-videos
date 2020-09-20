package sqlboiler

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud/service"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

func (r Repository) UpdateVideo(title string, videoDTO service.VideoDTO) (uuid.UUID, error) {
	video, err := r.FetchVideo(title)
	if err != nil {
		return uuid.UUID{}, err
	}
	tx, err := boil.BeginTx(r.ctx, nil)
	if err != nil {
		return uuid.UUID{}, err
	}
	if err := r.setCategoriesInVideo(videoDTO.Categories, video, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return uuid.UUID{}, err
		}
		return uuid.UUID{}, err
	}
	if err := r.setGenresInVideo(videoDTO.Genres, video, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return uuid.UUID{}, err
		}
		return uuid.UUID{}, err
	}
	video.Title = videoDTO.Title
	video.Description = videoDTO.Description
	video.YearLaunched = *videoDTO.YearLaunched
	video.Opened = null.Bool{Bool: videoDTO.Opened, Valid: true}
	video.Rating = int16(*videoDTO.Rating)
	video.Duration = *videoDTO.Duration
	var videoFile multipart.File
	fileName := null.String{}
	if videoDTO.File == nil {
		//TODO remove current video
	} else if videoDTO.File.Size > 0 {
		hash := sha256.New()
		var err error
		if videoFile, err = videoDTO.File.Open(); err != nil {
			return uuid.UUID{}, fmt.Errorf("could not genarete hash videoFile: %w", err)
		}
		if _, err := io.Copy(hash, videoFile); err != nil {
			return uuid.UUID{}, fmt.Errorf("could not genarete hash videoFile: %w", err)
		}
		hashName := fmt.Sprintf("%x", hash.Sum(nil))
		fileName = null.String{String: hashName, Valid: true}
	}
	video.VideoFile = fileName
	_, err = video.Update(r.ctx, tx, boil.Infer())
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return uuid.UUID{}, err
		}
		return uuid.UUID{}, fmt.Errorf("%s %w", videoDTO.Title, logger.ErrAlreadyExists)
	}
	videoID, err := uuid.Parse(video.ID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("could not parse video.ID: %v", err)
	}
	if _, err := r.repoFiles.UpdateFileToVideo(videoID, fileName.String, videoFile); err != nil {
		if err := tx.Rollback(); err != nil {
			return uuid.UUID{}, err
		}
		return uuid.UUID{}, err
	}
	if err := tx.Commit(); err != nil {
		return uuid.UUID{}, err
	}
	return videoID, nil
}

func (r Repository) AddVideo(videoDTO service.VideoDTO) (uuid.UUID, error) {
	id := uuid.New()
	var videoFile multipart.File
	fileName := null.String{}
	if videoDTO.File != nil && videoDTO.File.Size > 0 {
		hash := sha256.New()
		var err error
		if videoFile, err = videoDTO.File.Open(); err != nil {
			return uuid.UUID{}, fmt.Errorf("could not genarete hash videoFile: %w", err)
		}
		if _, err := io.Copy(hash, videoFile); err != nil {
			return uuid.UUID{}, fmt.Errorf("could not genarete hash videoFile: %w", err)
		}
		hashName := fmt.Sprintf("%x", hash.Sum(nil))
		fileName = null.String{String: hashName, Valid: true}
	}
	video := models.Video{
		ID:           id.String(),
		Title:        videoDTO.Title,
		Description:  videoDTO.Description,
		YearLaunched: *videoDTO.YearLaunched,
		Opened:       null.Bool{Bool: videoDTO.Opened, Valid: true},
		Rating:       int16(*videoDTO.Rating),
		Duration:     *videoDTO.Duration,
		VideoFile:    fileName,
	}
	tx, err := boil.BeginTx(r.ctx, nil)
	if err != nil {
		return uuid.UUID{}, err
	}
	err = video.Insert(r.ctx, tx, boil.Infer())
	if err != nil {
		var e *pq.Error
		if err := tx.Rollback(); err != nil {
			return uuid.UUID{}, err
		}
		if errors.As(err, &e) {
			if e.Code.Name() == "unique_violation" {
				return uuid.UUID{}, fmt.Errorf("title '%s' %w", videoDTO.Title, logger.ErrAlreadyExists)
			} else {
				return uuid.UUID{}, fmt.Errorf("%s: %w", "method Repository.AddVideo(videoDTO)", err)
			}
		}
	}
	if err := r.setCategoriesInVideo(videoDTO.Categories, video, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return uuid.UUID{}, err
		}
		return uuid.UUID{}, err
	}
	if err := r.setGenresInVideo(videoDTO.Genres, video, tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return uuid.UUID{}, err
		}
		return uuid.UUID{}, err
	}
	if videoDTO.File != nil && videoDTO.File.Size > 0 {
		if err := r.repoFiles.SaveFileToVideo(id, fileName.String, videoFile); err != nil {
			if err := tx.Rollback(); err != nil {
				return uuid.UUID{}, err
			}
			return uuid.UUID{}, fmt.Errorf("could not save file to video: %v", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (r Repository) setCategoriesInVideo(categories []service.Category, video models.Video, tx *sql.Tx) error {
	if categories == nil || len(categories) == 0 {
		return fmt.Errorf("none category is %w", logger.ErrNotFound)
	}
	clause := "name=?"
	categoryNames := make([]interface{}, len(categories))
	for i, category := range categories {
		if i > 0 {
			clause = fmt.Sprintf("name=? OR %s", clause)
		}
		categoryNames[i] = category.Name
	}
	categorySlice, err := models.Categories(
		Where(clause, categoryNames...),
	).AllG(r.ctx)
	if err != nil {
		return err
	}
	if len(categorySlice) == 0 {
		return fmt.Errorf("none category is %w", logger.ErrNotFound)
	}
	if err := video.SetCategories(r.ctx, tx, false, categorySlice...); err != nil {
		return fmt.Errorf("insert a new slice of categories and assign them to the video: %s", err)
	}
	return nil
}

func (r Repository) setGenresInVideo(genres []service.Genre, video models.Video, tx *sql.Tx) error {
	if genres == nil || len(genres) == 0 {
		return fmt.Errorf("none genre is %w", logger.ErrNotFound)
	}
	clause := "name=?"
	genreNames := make([]interface{}, len(genres))
	for i, genre := range genres {
		if i > 0 {
			clause = fmt.Sprintf("name=? OR %s", clause)
		}
		genreNames[i] = genre.Name
	}
	genreSlice, err := models.Genres(
		Where(clause, genreNames...),
	).AllG(r.ctx)
	if err != nil {
		return err
	}
	if len(genreSlice) == 0 {
		return fmt.Errorf("none genre is %w", logger.ErrNotFound)
	}
	if err := video.SetGenres(r.ctx, tx, false, genreSlice...); err != nil {
		return fmt.Errorf("insert a new slice of genres and assign them to the video: %s", err)
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
	videos, err := models.Videos(
		Load(models.VideoRels.Categories),
		Load(models.VideoRels.Genres),
		Limit(limit),
	).AllG(r.ctx)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (r Repository) FetchVideo(title string) (models.Video, error) {
	videoSlice, err := models.Videos(
		Load(models.VideoRels.Categories),
		Load(models.VideoRels.Genres),
		models.VideoWhere.Title.EQ(title),
	).AllG(r.ctx)
	if err != nil {
		return models.Video{}, err
	}
	if len(videoSlice) == 0 {
		return models.Video{}, sql.ErrNoRows
	}
	return *videoSlice[0], nil
}
