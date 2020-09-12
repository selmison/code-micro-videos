package service_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud/domain"
	"github.com/selmison/code-micro-videos/pkg/crud/service"
	"github.com/selmison/code-micro-videos/pkg/crud/service/mock"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/pkg/storage/sqlboiler"
	"github.com/selmison/code-micro-videos/testdata"
)

var (
	fakeYearLaunched       = new(int16)
	fakeDuration           = new(int16)
	fakeRating             = new(domain.VideoRating)
	fakeNotValidatedRating = new(domain.VideoRating)
)

func TestAddVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	const (
		fakeOpened     = false
		fakeGenreIndex = 0
	)
	fakeTitle := faker.Name()
	fakeDesc := faker.Sentence()
	*fakeYearLaunched = 2020
	*fakeDuration = 90
	*fakeRating = domain.TwelveRating
	*fakeNotValidatedRating = 111
	fakeExistGenreDTO := testdata.FakeGenresDTO[fakeGenreIndex]
	fakeDoesNotExistGenre := domain.Genre{Name: faker.FirstName()}
	fakeExistCategoryDTO := testdata.FakeCategoriesDTO[0]
	fakeDoesNotExistCategory := domain.Category{Name: faker.FirstName(), Description: &fakeDesc}
	type fields struct {
		r sqlboiler.Repository
	}
	type returns struct {
		video models.Video
		id    uuid.UUID
		err   error
	}
	type args struct {
		dto domain.VideoValidatable
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    returns
		wantErr bool
	}{
		{
			name:    "When VideoDTO is not provided",
			args:    args{domain.VideoValidatable{}},
			want:    returns{err: fmt.Errorf("'Title' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When the Title in VideoDTO is blank",
			args: args{domain.VideoValidatable{
				Title:        "    ",
				YearLaunched: fakeYearLaunched,
				Opened:       fakeOpened,
				Rating:       fakeRating,
				Duration:     fakeDuration,
			}},
			want:    returns{err: fmt.Errorf("'Title' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When the YearLaunched in VideoDTO is blank",
			args: args{domain.VideoValidatable{
				Title:        fakeTitle,
				YearLaunched: nil,
				Opened:       fakeOpened,
				Rating:       fakeRating,
				Duration:     fakeDuration,
			}},
			want:    returns{err: fmt.Errorf("'YearLaunched' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When the Rating in VideoDTO is blank",
			args: args{domain.VideoValidatable{
				Title:        fakeTitle,
				YearLaunched: fakeYearLaunched,
				Opened:       fakeOpened,
				Rating:       nil,
				Duration:     fakeDuration,
			}},
			want:    returns{err: fmt.Errorf("'Rating' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When the Duration in VideoDTO is blank",
			args: args{domain.VideoValidatable{
				Title:        fakeTitle,
				YearLaunched: fakeYearLaunched,
				Opened:       fakeOpened,
				Rating:       fakeRating,
				Duration:     nil,
			}},
			want:    returns{err: fmt.Errorf("'Duration' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When VideoDTO is with wrong genres and genres",
			args: args{
				domain.VideoValidatable{
					Title:        fakeTitle,
					Description:  fakeDesc,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
					Genres:       []domain.Genre{fakeDoesNotExistGenre},
					Categories:   []domain.Category{fakeDoesNotExistCategory},
				},
			},
			want:    returns{err: logger.ErrIsRequired},
			wantErr: true,
		},
		{
			name: "When VideoDTO is without genres and genres",
			args: args{
				domain.VideoValidatable{
					Title:        fakeTitle,
					Description:  fakeDesc,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
				},
			},
			want:    returns{err: fmt.Errorf("'Categories' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When VideoDTO is right",
			args: args{domain.VideoValidatable{
				Title:        fakeTitle,
				Description:  fakeDesc,
				YearLaunched: fakeYearLaunched,
				Opened:       fakeOpened,
				Rating:       fakeRating,
				Duration:     fakeDuration,
				Genres:       []domain.Genre{fakeExistGenreDTO},
				Categories:   []domain.Category{fakeExistCategoryDTO},
			}},
			want: returns{video: models.Video{
				Title:        fakeTitle,
				Description:  fakeDesc,
				Opened:       null.BoolFrom(fakeOpened),
				YearLaunched: *fakeYearLaunched,
				Rating:       int16(*fakeRating),
				Duration:     *fakeDuration,
			},
				err: nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr || tt.name == "When VideoDTO is with wrong genres and genres" {
				tt.args.dto.Title = strings.ToLower(strings.TrimSpace(tt.args.dto.Title))
				mockR.EXPECT().
					AddVideo(tt.args.dto).
					Return(tt.want.id, tt.want.err)
			}
			s := service.NewService(mockR)
			_, err := s.CreateVideo(tt.args.ctx, tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddVideo() error = '%v', wantErr '%v'", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("AddVideo() got: \"%v\", want: \"%v\"", err, tt.want.err)
			}
		})
	}
}

func Test_service_RemoveVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeExistTitle := testdata.FakeVideos[0].Title
	const fakeDoesNotExistTitle = "fakeDoesNotExistTitle"
	type args struct {
		title string
	}
	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name:    "When title is blank",
			args:    args{"     "},
			want:    fmt.Errorf("'title' %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name:    "When title is not found",
			args:    args{fakeDoesNotExistTitle},
			want:    fmt.Errorf("%s: %w", fakeDoesNotExistTitle, logger.ErrNotFound),
			wantErr: true,
		},
		{
			name:    "When title is found",
			args:    args{fakeExistTitle},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "When title is not found" {
				tt.args.title = strings.ToLower(strings.ToLower(tt.args.title))
				mockR.EXPECT().
					RemoveVideo(tt.args.title).
					Return(tt.want)
			} else if tt.name == "When title is found" {
				mockR.EXPECT().
					RemoveVideo(tt.args.title).
					Return(tt.want)
			}
			s := service.NewService(mockR)
			err := s.RemoveVideo(tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveVideo() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("RemoveVideo() got: '%v', want: '%v'", err, tt.want)
			}
		})
	}
}

func Test_service_UpdateVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	const (
		fakeExistTitle        = "fakeExistTitle"
		fakeDoesNotExistTitle = "fakeDoesNotExistTitle"
		fakeOpened            = false
		fakeCategoryIndex     = 0
		fakeGenreIndex        = 0
	)
	fakeTitle := faker.Name()
	fakeDesc := faker.Sentence()
	*fakeYearLaunched = 2020
	*fakeDuration = 90
	*fakeRating = domain.TwelveRating
	*fakeNotValidatedRating = 111
	fakeExistGenreDTO := testdata.FakeGenresDTO[fakeGenreIndex]
	fakeDoesNotExistGenre := domain.Genre{Name: faker.FirstName()}
	fakeExistCategoryDTO := testdata.FakeCategoriesDTO[fakeCategoryIndex]
	fakeDoesNotExistCategory := domain.Category{Name: faker.FirstName(), Description: &fakeDesc}

	type fields struct {
		r sqlboiler.Repository
	}
	type args struct {
		title string
		dto   domain.VideoValidatable
	}
	type returns struct {
		id  uuid.UUID
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "When title is blank",
			args: args{
				"     ",
				domain.VideoValidatable{
					Title: faker.Name(),
				},
			},
			want:    returns{err: fmt.Errorf("'title' %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name:    "When VideoDTO is not provided",
			args:    args{fakeExistTitle, domain.VideoValidatable{}},
			want:    returns{err: fmt.Errorf("'Title' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When the Title in VideoDTO is blank",
			args: args{
				title: fakeExistTitle,
				dto: domain.VideoValidatable{
					Title:        "    ",
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
				}},
			want:    returns{err: fmt.Errorf("'Title' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When the YearLaunched in VideoDTO is blank",
			args: args{
				title: fakeExistTitle,
				dto: domain.VideoValidatable{
					Title:        fakeTitle,
					YearLaunched: nil,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
				}},
			want:    returns{err: fmt.Errorf("'YearLaunched' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When the Rating in VideoDTO is blank",
			args: args{
				title: fakeTitle,
				dto: domain.VideoValidatable{
					Title:        fakeTitle,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       nil,
					Duration:     fakeDuration,
				}},
			want:    returns{err: fmt.Errorf("'Rating' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When the Duration in VideoDTO is blank",
			args: args{
				title: fakeTitle,
				dto: domain.VideoValidatable{
					Title:        fakeTitle,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     nil,
				},
			},
			want:    returns{err: fmt.Errorf("'Duration' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When VideoDTO is with wrong genres and genres",
			args: args{
				fakeExistTitle,
				domain.VideoValidatable{
					Title:        fakeTitle,
					Description:  fakeDesc,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
					Genres:       []domain.Genre{fakeDoesNotExistGenre},
					Categories:   []domain.Category{fakeDoesNotExistCategory},
				},
			},
			want:    returns{err: logger.ErrIsRequired},
			wantErr: true,
		},
		{
			name: "When Category is without genres and genres",
			args: args{
				fakeExistTitle,
				domain.VideoValidatable{
					Title:        fakeTitle,
					Description:  fakeDesc,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
				},
			},
			want:    returns{err: fmt.Errorf("'Categories' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When title is not found",
			args: args{
				fakeDoesNotExistTitle,
				domain.VideoValidatable{
					Title:        fakeTitle,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
					Genres:       []domain.Genre{fakeExistGenreDTO},
					Categories:   []domain.Category{fakeExistCategoryDTO},
				},
			},
			want:    returns{err: fmt.Errorf("%s: %w", fakeDoesNotExistTitle, logger.ErrNotFound)},
			wantErr: true,
		},
		{
			name: "When title is found and VideoDTO is right",
			args: args{
				fakeExistTitle,
				domain.VideoValidatable{
					Title:        fakeTitle,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
					Genres:       []domain.Genre{fakeExistGenreDTO},
					Categories:   []domain.Category{fakeExistCategoryDTO},
				},
			},
			want:    returns{err: nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "When the Title in Category already exists" ||
				tt.name == "When VideoDTO is with wrong genres and genres" ||
				tt.name == "When title is not found" ||
				tt.name == "When title is found and VideoDTO is right" {
				dto := domain.VideoValidatable{
					Title:        strings.ToLower(strings.TrimSpace(tt.args.dto.Title)),
					Description:  tt.args.dto.Description,
					YearLaunched: tt.args.dto.YearLaunched,
					Opened:       tt.args.dto.Opened,
					Rating:       tt.args.dto.Rating,
					Duration:     tt.args.dto.Duration,
					Categories:   tt.args.dto.Categories,
					Genres:       tt.args.dto.Genres,
				}
				title := strings.ToLower(strings.TrimSpace(tt.args.title))
				mockR.EXPECT().
					UpdateVideo(title, dto).
					Return(tt.want.id, tt.want.err)
			}
			s := service.NewService(mockR)
			_, err := s.UpdateVideo(tt.args.title, tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateVideo() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("AddVideo() got: \"%v\", want: \"%v\"", err, tt.want)
			}
		})
	}
}

func Test_service_GetVideos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeVideoSlice := testdata.FakeVideoSlice
	const fakeLimit = testdata.FakeVideosLength
	type args struct {
		limit int
	}
	type returns struct {
		videos models.VideoSlice
		err    error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name:    "When limit is less than zero",
			args:    args{-1},
			want:    returns{nil, logger.ErrInvalidedLimit},
			wantErr: true,
		},
		{
			name:    "When limit is right",
			args:    args{fakeLimit},
			want:    returns{fakeVideoSlice, nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockR.EXPECT().
					GetVideos(tt.args.limit).
					Return(
						fakeVideoSlice,
						nil,
					)
			}
			s := service.NewService(mockR)
			got, err := s.GetVideos(tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVideos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.videos) {
				t.Errorf("GetVideos() got = %v, want %v", got, tt.want.videos)
			}
			if !reflect.DeepEqual(err, tt.want.err) {
				t.Errorf("GetVideos() got = %v, want %v", err, tt.want.err)
			}
		})
	}
}

func Test_service_FetchVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	const (
		fakeDoesNotExistTitle = "fakeDoesNotExistTitle"
		fakeExistTitle        = "fakeExistTitle"
	)
	fakeErrorInternalApplication := fmt.Errorf("anyname: %w", logger.ErrInternalApplication)
	fakeVideo := models.Video{
		ID:           uuid.New().String(),
		Title:        fakeExistTitle,
		Description:  faker.Sentence(),
		YearLaunched: 2020,
		Opened:       null.Bool{Bool: true, Valid: true},
		Rating:       int16(domain.TwelveRating),
		Duration:     150,
	}
	type args struct {
		title string
	}
	type returns struct {
		video models.Video
		err   error
	}
	tests := []struct {
		name       string
		args       args
		want       returns
		wantErr    bool
		setupMockR func()
	}{
		{
			name: "When throw the error internal application",
			args: args{"anyName"},
			want: returns{
				models.Video{},
				fakeErrorInternalApplication,
			},
			wantErr: true,
			setupMockR: func() {
				mockR.EXPECT().
					FetchVideo("anyname").
					Return(
						models.Video{},
						fakeErrorInternalApplication,
					)
			},
		},
		{
			name: "When title is not found",
			args: args{fakeDoesNotExistTitle},
			want: returns{
				models.Video{},
				fmt.Errorf("%s: %w", strings.ToLower(strings.ToLower(fakeDoesNotExistTitle)), logger.ErrNotFound),
			},
			wantErr: true,
			setupMockR: func() {
				mockR.EXPECT().
					FetchVideo(strings.ToLower(strings.ToLower(fakeDoesNotExistTitle))).
					Return(
						models.Video{},
						sql.ErrNoRows,
					)
			},
		},
		{
			name: "When title is found",
			args: args{fakeExistTitle},
			want: returns{
				fakeVideo,
				nil,
			},
			wantErr: false,
			setupMockR: func() {
				mockR.EXPECT().
					FetchVideo(strings.ToLower(strings.TrimSpace(fakeExistTitle))).
					Return(
						fakeVideo,
						nil,
					)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMockR()
			s := service.NewService(mockR)
			got, err := s.FetchVideo(tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchVideo() error: %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got, tt.want.video, "they should be equal")
			if err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("GetVideo() got: %v, want: %v", err, tt.want.err)
			}
		})
	}
}
