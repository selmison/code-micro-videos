package crud_test

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
	"github.com/selmison/code-micro-videos/pkg/crud"
	"github.com/selmison/code-micro-videos/pkg/crud/mock"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/pkg/storage/sqlboiler"
	"github.com/selmison/code-micro-videos/testdata"
)

var (
	fakeYearLaunched       = new(int16)
	fakeDuration           = new(int16)
	fakeRating             = new(crud.VideoRating)
	fakeNotValidatedRating = new(crud.VideoRating)
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
	*fakeRating = crud.TwelveRating
	*fakeNotValidatedRating = 111
	fakeExistGenreDTO := testdata.FakeGenresDTO[fakeGenreIndex]
	fakeDoesNotExistGenre := crud.GenreDTO{Name: faker.FirstName()}
	fakeExistCategoryDTO := testdata.FakeCategoriesDTO[0]
	fakeDoesNotExistCategory := crud.CategoryDTO{Name: faker.FirstName(), Description: faker.Sentence()}
	type fields struct {
		r sqlboiler.Repository
	}
	type returns struct {
		video models.Video
		err   error
	}
	type args struct {
		dto crud.VideoDTO
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
			args:    args{crud.VideoDTO{}},
			want:    returns{err: fmt.Errorf("'Title' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When the Title in VideoDTO is blank",
			args: args{crud.VideoDTO{
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
			args: args{crud.VideoDTO{
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
			args: args{crud.VideoDTO{
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
			args: args{crud.VideoDTO{
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
			name: "When VideoDTO is with wrong categories and genres",
			args: args{
				crud.VideoDTO{
					Title:        fakeTitle,
					Description:  fakeDesc,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
					Genres:       []crud.GenreDTO{fakeDoesNotExistGenre},
					Categories:   []crud.CategoryDTO{fakeDoesNotExistCategory},
				},
			},
			want:    returns{err: logger.ErrIsRequired},
			wantErr: true,
		},
		{
			name: "When VideoDTO is without categories and genres",
			args: args{
				crud.VideoDTO{
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
			args: args{crud.VideoDTO{
				Title:        fakeTitle,
				Description:  fakeDesc,
				YearLaunched: fakeYearLaunched,
				Opened:       fakeOpened,
				Rating:       fakeRating,
				Duration:     fakeDuration,
				Genres:       []crud.GenreDTO{fakeExistGenreDTO},
				Categories:   []crud.CategoryDTO{fakeExistCategoryDTO},
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
			if !tt.wantErr || tt.name == "When VideoDTO is with wrong categories and genres" {
				tt.args.dto.Title = strings.ToLower(strings.TrimSpace(tt.args.dto.Title))
				mockR.EXPECT().
					AddVideo(tt.args.dto).
					Return(tt.want.err)
			}
			s := crud.NewService(mockR)
			err := s.AddVideo(tt.args.dto)
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
			s := crud.NewService(mockR)
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
	fakeDescription := faker.Sentence()
	*fakeYearLaunched = 2020
	*fakeDuration = 90
	*fakeRating = crud.TwelveRating
	*fakeNotValidatedRating = 111
	fakeExistGenreDTO := testdata.FakeGenresDTO[fakeGenreIndex]
	fakeDoesNotExistGenre := crud.GenreDTO{Name: faker.FirstName()}
	fakeExistCategoryDTO := testdata.FakeCategoriesDTO[fakeCategoryIndex]
	fakeDoesNotExistCategory := crud.CategoryDTO{Name: faker.FirstName(), Description: faker.Sentence()}

	type fields struct {
		r sqlboiler.Repository
	}
	type args struct {
		title string
		dto   crud.VideoDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "When title is blank",
			args: args{
				"     ",
				crud.VideoDTO{
					Title: faker.Name(),
				},
			},
			want:    fmt.Errorf("'title' %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name:    "When VideoDTO is not provided",
			args:    args{fakeExistTitle, crud.VideoDTO{}},
			want:    fmt.Errorf("'Title' field %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name: "When the Title in VideoDTO is blank",
			args: args{
				title: fakeExistTitle,
				dto: crud.VideoDTO{
					Title:        "    ",
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
				}},
			want:    fmt.Errorf("'Title' field %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name: "When the YearLaunched in VideoDTO is blank",
			args: args{
				title: fakeExistTitle,
				dto: crud.VideoDTO{
					Title:        fakeTitle,
					YearLaunched: nil,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
				}},
			want:    fmt.Errorf("'YearLaunched' field %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name: "When the Rating in VideoDTO is blank",
			args: args{
				title: fakeTitle,
				dto: crud.VideoDTO{
					Title:        fakeTitle,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       nil,
					Duration:     fakeDuration,
				}},
			want:    fmt.Errorf("'Rating' field %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name: "When the Duration in VideoDTO is blank",
			args: args{
				title: fakeTitle,
				dto: crud.VideoDTO{
					Title:        fakeTitle,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     nil,
				},
			},
			want:    fmt.Errorf("'Duration' field %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name: "When VideoDTO is with wrong categories and genres",
			args: args{
				fakeExistTitle,
				crud.VideoDTO{
					Title:        fakeTitle,
					Description:  fakeDescription,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
					Genres:       []crud.GenreDTO{fakeDoesNotExistGenre},
					Categories:   []crud.CategoryDTO{fakeDoesNotExistCategory},
				},
			},
			want:    logger.ErrIsRequired,
			wantErr: true,
		},
		{
			name: "When CategoryDTO is without categories and genres",
			args: args{
				fakeExistTitle,
				crud.VideoDTO{
					Title:        fakeTitle,
					Description:  fakeDescription,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
				},
			},
			want:    fmt.Errorf("'Categories' field %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name: "When title is not found",
			args: args{
				fakeDoesNotExistTitle,
				crud.VideoDTO{
					Title:        fakeTitle,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
					Genres:       []crud.GenreDTO{fakeExistGenreDTO},
					Categories:   []crud.CategoryDTO{fakeExistCategoryDTO},
				},
			},
			want:    fmt.Errorf("%s: %w", fakeDoesNotExistTitle, logger.ErrNotFound),
			wantErr: true,
		},
		{
			name: "When title is found and VideoDTO is right",
			args: args{
				fakeExistTitle,
				crud.VideoDTO{
					Title:        fakeTitle,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
					Genres:       []crud.GenreDTO{fakeExistGenreDTO},
					Categories:   []crud.CategoryDTO{fakeExistCategoryDTO},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "When the Title in CategoryDTO already exists" ||
				tt.name == "When VideoDTO is with wrong categories and genres" ||
				tt.name == "When title is not found" ||
				tt.name == "When title is found and VideoDTO is right" {
				dto := crud.VideoDTO{
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
					Return(tt.want)
			}
			s := crud.NewService(mockR)
			err := s.UpdateVideo(tt.args.title, tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateVideo() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.want.Error() {
				t.Errorf("AddVideo() got: \"%v\", want: \"%v\"", err, tt.want)
			}

			//if !errors.Is(err, tt.want) {
			//	t.Errorf("UpdateVideo() got: '%v', want: '%v'", err, tt.want)
			//}
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
			s := crud.NewService(mockR)
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
		Rating:       int16(crud.TwelveRating),
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
			s := crud.NewService(mockR)
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
