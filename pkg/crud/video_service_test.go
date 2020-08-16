package crud_test

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud"
	"github.com/selmison/code-micro-videos/pkg/crud/mock"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/pkg/storage/sqlboiler"
	"github.com/selmison/code-micro-videos/testdata"
)

func TestAddVideo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeTitle := faker.Name()
	fakeDesc := faker.Sentence()
	const (
		fakeYearLaunched int16 = 2020
		fakeOpened             = false
		fakeRating             = crud.TwelveRating
		fakeDuration     int16 = 90
	)
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
			want:    returns{models.Video{}, logger.ErrIsRequired},
			wantErr: true,
		},
		{
			name: "When the title in VideoDTO is blank",
			args: args{crud.VideoDTO{
				Title: "    ",
			}},
			want:    returns{models.Video{}, logger.ErrIsRequired},
			wantErr: true,
		},
		{
			name: "When video rating in VideoDTO is not validated",
			args: args{crud.VideoDTO{
				Title:  fakeTitle,
				Rating: 111,
			}},
			want:    returns{models.Video{}, logger.ErrIsNotValidated},
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
			}},
			want: returns{models.Video{
				Title:        fakeTitle,
				Description:  fakeDesc,
				YearLaunched: fakeYearLaunched,
				Opened:       null.BoolFrom(fakeOpened),
				Rating:       int16(fakeRating),
				Duration:     fakeDuration,
			}, nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
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
			if !errors.Is(err, tt.want.err) {
				t.Errorf("AddVideo() got = '%v', want '%v'", err, tt.want.err)
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
			want:    fmt.Errorf("'name' %w", logger.ErrIsRequired),
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
	)
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
			want:    logger.ErrIsRequired,
			wantErr: true,
		},
		{
			name:    "When VideoDTO is not provided",
			args:    args{fakeExistTitle, crud.VideoDTO{}},
			want:    logger.ErrIsRequired,
			wantErr: true,
		},
		{
			name: "When title is not found",
			args: args{
				fakeDoesNotExistTitle,
				crud.VideoDTO{
					Title: faker.FirstName(),
				},
			},
			want:    fmt.Errorf("%s: %w", fakeDoesNotExistTitle, logger.ErrNotFound),
			wantErr: true,
		},
		{
			name: "When title is found and VideoDTO is provided",
			args: args{
				fakeExistTitle,
				crud.VideoDTO{
					Title: faker.FirstName(),
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "When title is not found" || tt.name == "When title is found and VideoDTO is provided" {
				mockR.EXPECT().
					UpdateVideo(tt.args.title, tt.args.dto).
					Return(tt.want)
			}
			s := crud.NewService(mockR)
			err := s.UpdateVideo(tt.args.title, tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateVideo() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !errors.Is(err, tt.want) {
				t.Errorf("UpdateVideo() got: '%v', want: '%v'", err, tt.want)
			}
		})
	}
}

func Test_service_GetVideos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeVideoSlice := testdata.FakeVideoSlice
	fakeLimit := len(fakeVideoSlice)
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
	fakeErrorInternalApplication := fmt.Errorf("Service.FetchVideo(): %w", logger.ErrInternalApplication)
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
		name string
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
					FetchVideo("anyName").
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
				fmt.Errorf("%s: %w", fakeDoesNotExistTitle, logger.ErrNotFound),
			},
			wantErr: true,
			setupMockR: func() {
				mockR.EXPECT().
					FetchVideo(fakeDoesNotExistTitle).
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
					FetchVideo(fakeExistTitle).
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
			got, err := s.FetchVideo(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchVideo() error: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.video) {
				t.Errorf("GetVideo() got: %v, want: %v", got, tt.want.video)
			}
			if tt.wantErr && errors.Is(err, tt.want.err) {
				t.Errorf("GetVideo() got: %v, want: %v", err, tt.want.err)
			}
		})
	}
}
