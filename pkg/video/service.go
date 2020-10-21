package video

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/go-kit/kit/log"

	"github.com/selmison/code-micro-videos/pkg/id_generator"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/pkg/validator"
)

// Service manages videos.
type Service interface {
	// Create creates a new video.
	Create(ctx context.Context, newVideo NewVideo) (Video, error)

	// Destroy destroys a video.
	Destroy(ctx context.Context, id string) error

	// List returns a list of videos.
	List(ctx context.Context) ([]Video, error)

	// Show returns the details of a video.
	Show(ctx context.Context, id string) (Video, error)

	// Update updates an existing video.
	Update(ctx context.Context, id string, updateVideo UpdateVideo) error
}

type VideoRating int16

const (
	FreeRating VideoRating = iota + 1
	TenRating
	TwelveRating
	FourteenRating
	SixteenRating
	EighteenRating
)

func (v *VideoRating) String() string {
	return [...]string{"Free", "10", "12", "14", "16", "18"}[*v-1]
}

func (v *VideoRating) Validate() error {
	switch *v {
	case FreeRating, TenRating, TwelveRating, FourteenRating, SixteenRating, EighteenRating:
		return nil
	}
	return fmt.Errorf("video rating %w", logger.ErrIsNotValidated)
}

// Video represents a single Video.
type Video struct {
	Id               string
	Title            string
	Description      string
	YearLaunched     int16
	Opened           bool
	Rating           VideoRating
	Duration         int16
	CategoriesId     []string
	GenresId         []string
	VideoFileHandler *multipart.FileHeader
}

// NewVideo contains the details of a new Video.
type NewVideo struct {
	Title            string
	Description      string
	YearLaunched     *int16
	Opened           bool
	Rating           *VideoRating
	Duration         *int16
	CategoriesId     []string
	GenresId         []string
	VideoFileHandler *multipart.FileHeader
}

func (v *NewVideo) Validate() error {
	errField := " the %s field"
	v.Title = strings.TrimSpace(v.Title)
	if err := validator.CheckAll(
		validator.NewCheck(
			v.Title == "",
			&logger.ResultError{
				ErrMsg: fmt.Sprintf(errField, "Title"),
				Err:    logger.ErrIsRequired,
			},
		),
		validator.NewCheck(
			v.YearLaunched == nil,
			&logger.ResultError{
				ErrMsg: fmt.Sprintf(errField, "YearLaunched"),
				Err:    logger.ErrIsRequired,
			}),
		validator.NewCheck(
			v.Rating == nil,
			&logger.ResultError{
				ErrMsg: fmt.Sprintf(errField, "Rating"),
				Err:    logger.ErrIsRequired,
			},
		),
		validator.NewCheck(
			v.Duration == nil,
			&logger.ResultError{
				ErrMsg: fmt.Sprintf(errField, "Duration"),
				Err:    logger.ErrIsRequired,
			},
		),
		validator.NewCheck(
			v.Duration != nil && *v.Duration < 0,
			&logger.ResultError{
				ErrMsg: fmt.Sprintf(errField, "Duration"),
				Err:    logger.ErrIsNotValidated,
			},
		),
		validator.NewCheck(
			len(v.CategoriesId) == 0,
			&logger.ResultError{
				ErrMsg: fmt.Sprintf(errField, "CategoriesId"),
				Err:    logger.ErrIsRequired,
			},
		),
		validator.NewCheck(
			len(v.GenresId) == 0,
			&logger.ResultError{
				ErrMsg: fmt.Sprintf(errField, "GenresId"),
				Err:    logger.ErrIsRequired,
			},
		),
	); err != nil {
		return err
	}
	if err := v.Rating.Validate(); err != nil {
		return err
	}
	return nil
}

func (v NewVideo) ToVideo(id string) Video {
	if v.YearLaunched == nil ||
		v.Rating == nil ||
		v.Duration == nil {
		return Video{}
	}
	return Video{
		Id:               id,
		Title:            v.Title,
		Description:      v.Description,
		YearLaunched:     *v.YearLaunched,
		Opened:           v.Opened,
		Rating:           *v.Rating,
		Duration:         *v.Duration,
		CategoriesId:     v.CategoriesId,
		GenresId:         v.GenresId,
		VideoFileHandler: v.VideoFileHandler,
	}
}

// UpdateVideo contains updates of an existing video.
type UpdateVideo struct {
	Title            *string
	Description      *string
	YearLaunched     *int16
	Opened           *bool
	Rating           *VideoRating
	Duration         *int16
	CategoriesId     []string
	GenresId         []string
	VideoFileHandler *multipart.FileHeader
}

func (v UpdateVideo) update(video Video) Video {
	if v.Title != nil {
		video.Title = *v.Title
	}
	if v.Description != nil {
		video.Description = *v.Description
	}
	if v.YearLaunched != nil {
		video.YearLaunched = *v.YearLaunched
	}
	if v.Opened != nil {
		video.Opened = *v.Opened
	}
	if v.Rating != nil {
		video.Rating = *v.Rating
	}
	if v.Duration != nil {
		video.Duration = *v.Duration
	}
	if v.CategoriesId != nil {
		video.CategoriesId = v.CategoriesId
	}
	if v.GenresId != nil {
		video.GenresId = v.GenresId
	}
	if v.VideoFileHandler != nil {
		video.VideoFileHandler = v.VideoFileHandler
	}
	return video
}

func (v *UpdateVideo) Validate() error {
	if v.isEmpty() {
		return fmt.Errorf("'%s' param %w", "updateVideo", logger.ErrCouldNotBeEmpty)
	}
	errField := "the %s field"
	if v.Title != nil {
		*v.Title = strings.TrimSpace(*v.Title)
	}
	if err := validator.CheckAll(
		validator.NewCheck(
			v.Title != nil && *v.Title == "",
			&logger.ResultError{
				ErrMsg: fmt.Sprintf(errField, "Title"),
				Err:    logger.ErrCouldNotBeEmpty,
			},
		),
		validator.NewCheck(
			v.YearLaunched != nil && *v.YearLaunched < 0,
			&logger.ResultError{
				ErrMsg: fmt.Sprintf(errField, "YearLaunched"),
				Err:    logger.ErrIsNotValidated,
			},
		),
		validator.NewCheck(
			v.Duration != nil && *v.Duration < 0,
			&logger.ResultError{
				ErrMsg: fmt.Sprintf(errField, "Duration"),
				Err:    logger.ErrIsNotValidated,
			},
		),
	); err != nil {
		return err
	}
	if v.Rating != nil {
		if err := v.Rating.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type service struct {
	idGenerator id_generator.IdGenerator
	repo        Repository
}

// NewService returns a new Service with all of the expected middlewares wired in.
func NewService(idGenerator id_generator.IdGenerator, repo Repository, logger log.Logger) Service {
	var svc Service
	{
		svc = service{idGenerator: idGenerator, repo: repo}
		svc = NewValidationMiddleware()(svc)
		svc = NewLoggingMiddleware(logger)(svc)
	}
	return svc
}

func (svc service) Create(ctx context.Context, newVideo NewVideo) (Video, error) {
	id, err := svc.idGenerator.Generate()
	if err != nil {
		return Video{}, err
	}
	video := newVideo.ToVideo(id)
	err = svc.repo.Store(ctx, video)
	if err != nil {
		return Video{}, err
	}
	return video, nil
}

func (svc service) Destroy(ctx context.Context, id string) error {
	return svc.repo.DeleteOne(ctx, id)
}

func (svc service) List(ctx context.Context) ([]Video, error) {
	videos, err := svc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func (svc service) Show(ctx context.Context, id string) (Video, error) {
	video, err := svc.repo.GetOne(ctx, id)
	if err != nil {
		return Video{}, err
	}
	return video, nil
}

func (svc service) Update(ctx context.Context, id string, updateVideo UpdateVideo) error {
	if updateVideo.Title != nil {
		fmt.Println("Title-2: ", *updateVideo.Title)
	}
	fmt.Printf("2: %v %v %#v\n", ctx, id, updateVideo)
	return svc.repo.UpdateOne(ctx, id, updateVideo)
}

func (v NewVideo) isEmpty() bool {
	return v.compare(NewVideo{})
}

func (v NewVideo) compare(b NewVideo) bool {
	if &v == &b {
		return true
	}
	if v.Title != b.Title {
		return false
	}
	if v.Description != b.Description {
		return false
	}
	if v.YearLaunched == nil && b.YearLaunched != nil || v.YearLaunched != nil && b.YearLaunched == nil {
		return false
	}

	if v.YearLaunched != nil && b.YearLaunched != nil && *v.YearLaunched != *b.YearLaunched {
		return false
	}
	if v.Opened != b.Opened {
		return false
	}

	if v.Rating == nil && b.Rating != nil || v.Rating != nil && b.Rating == nil {
		return false
	}
	if v.Rating != nil && b.Rating != nil && *v.Rating != *b.Rating {
		return false
	}

	if v.Duration == nil && b.Duration != nil || v.Duration != nil && b.Duration == nil {
		return false
	}
	if v.Duration != nil && b.Duration != nil && *v.Duration != *b.Duration {
		return false
	}

	if len(v.GenresId) != len(b.GenresId) {
		return false
	}
	for i, v := range v.GenresId {
		if b.GenresId[i] != v {
			return false
		}
	}

	if len(v.CategoriesId) != len(b.CategoriesId) {
		return false
	}
	for i, v := range v.CategoriesId {
		if b.CategoriesId[i] != v {
			return false
		}
	}
	if v.VideoFileHandler == nil && b.VideoFileHandler != nil || v.VideoFileHandler != nil && b.VideoFileHandler == nil {
		return false
	}
	//TODO: add deep compare to VideoFileHandler
	return true
}

func (v UpdateVideo) isEmpty() bool {
	return v.compare(UpdateVideo{})
}

func (v UpdateVideo) compare(b UpdateVideo) bool {
	if &v == &b {
		return true
	}

	if v.Title == nil && b.Title != nil || v.Title != nil && b.Title == nil {
		return false
	}
	if v.Title != nil && b.Title != nil && *v.Title != *b.Title {
		return false
	}

	if v.Description == nil && b.Description != nil || v.Description != nil && b.Description == nil {
		return false
	}
	if v.Description != nil && b.Description != nil && *v.Description != *b.Description {
		return false
	}

	if v.YearLaunched == nil && b.YearLaunched != nil || v.YearLaunched != nil && b.YearLaunched == nil {
		return false
	}
	if v.YearLaunched != nil && b.YearLaunched != nil && *v.YearLaunched != *b.YearLaunched {
		return false
	}

	if v.Opened == nil && b.Opened != nil || v.Opened != nil && b.Opened == nil {
		return false
	}
	if v.Opened != nil && b.Opened != nil && *v.Opened != *b.Opened {
		return false
	}

	if v.Rating == nil && b.Rating != nil || v.Rating != nil && b.Rating == nil {
		return false
	}
	if v.Rating != nil && b.Rating != nil && *v.Rating != *b.Rating {
		return false
	}

	if v.Duration == nil && b.Duration != nil || v.Duration != nil && b.Duration == nil {
		return false
	}
	if v.Duration != nil && b.Duration != nil && *v.Duration != *b.Duration {
		return false
	}

	if len(v.GenresId) != len(b.GenresId) {
		return false
	}
	for i, v := range v.GenresId {
		if b.GenresId[i] != v {
			return false
		}
	}

	if len(v.CategoriesId) != len(b.CategoriesId) {
		return false
	}
	for i, v := range v.CategoriesId {
		if b.CategoriesId[i] != v {
			return false
		}
	}

	if v.VideoFileHandler == nil && b.VideoFileHandler != nil || v.VideoFileHandler != nil && b.VideoFileHandler == nil {
		return false
	}
	//TODO: add deep compare to VideoFileHandler
	return true
}
