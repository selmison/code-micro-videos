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
	tx, err := boil.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}
	_, err = video.UpdateG(r.ctx, boil.Infer())
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("%s %w", videoDTO.Title, logger.ErrAlreadyExists)
	}
	if err := r.setCategoriesInVideo(videoDTO.Categories, video, tx); err != nil {
		return err
	}
	if err := r.setGenresInVideo(videoDTO.Genres, video, tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r Repository) AddVideo(videoDTO crud.VideoDTO) error {
	video := models.Video{
		ID:    uuid.New().String(),
		Title: videoDTO.Title,
	}
	tx, err := boil.BeginTx(r.ctx, nil)
	if err != nil {
		return err
	}
	err = video.InsertG(r.ctx, boil.Infer())
	if err != nil {
		var e *pq.Error
		if err := tx.Rollback(); err != nil {
			return err
		}
		if errors.As(err, &e) {
			if e.Code.Name() == "unique_violation" {
				return fmt.Errorf("title '%s' %w", videoDTO.Title, logger.ErrAlreadyExists)
			} else {
				return fmt.Errorf("%s: %w", "method Repository.AddVideo(videoDTO)", err)
			}
		}
	}
	if err := r.setCategoriesInVideo(videoDTO.Categories, video, tx); err != nil {
		return err
	}
	if err := r.setGenresInVideo(videoDTO.Genres, video, tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r Repository) setCategoriesInVideo(categories []crud.CategoryDTO, video models.Video, tx *sql.Tx) error {
	if categories != nil {
		clause := "name=?"
		categoryNames := make([]interface{}, len(categories))
		for i, category := range categories {
			if i > 0 {
				clause = fmt.Sprintf("name=? OR %s", clause)
			}
			categoryNames[i] = strings.ToLower(category.Name)
		}
		categorySlice, err := models.Categories(
			Where(clause, categoryNames...),
		).AllG(r.ctx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
		if len(categorySlice) == 0 {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("none category is %w", logger.ErrNotFound)
		}
		if err := video.SetCategoriesG(r.ctx, false, categorySlice...); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("insert a new slice of categories and assign them to the video: %s", err)
		}
	}
	return nil
}

func (r Repository) setGenresInVideo(genres []crud.GenreDTO, video models.Video, tx *sql.Tx) error {
	if genres != nil {
		clause := "name=?"
		genreNames := make([]interface{}, len(genres))
		for i, genre := range genres {
			if i > 0 {
				clause = fmt.Sprintf("name=? OR %s", clause)
			}
			genreNames[i] = strings.ToLower(genre.Name)
		}
		genreSlice, err := models.Genres(
			Where(clause, genreNames...),
		).AllG(r.ctx)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
		if len(genreSlice) == 0 {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("none genre is %w", logger.ErrNotFound)
		}
		if err := video.SetGenresG(r.ctx, false, genreSlice...); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("insert a new slice of genres and assign them to the video: %s", err)
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
