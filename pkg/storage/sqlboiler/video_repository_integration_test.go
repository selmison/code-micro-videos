// +build integration

package sqlboiler

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/testdata"
)

func TestRepository_AddVideo(t *testing.T) {
	cfg, teardownTestCase, repository, err := setupTestCase(nil)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	const (
		fakeExistTitle        = "action"
		fakeDoesNotExistTitle = "fakeDoesNotExistTitle"
	)
	fakeExistVideo := models.Video{
		ID:    uuid.New().String(),
		Title: fakeExistTitle,
	}
	type args struct {
		videoDTO crud.VideoDTO
	}
	type returns struct {
		video models.Video
		err   error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "When title in VideoDTO already exists",
			args: args{crud.VideoDTO{
				Title: fakeExistTitle,
			}},
			want: returns{
				models.Video{
					Title: fakeExistTitle,
				},
				nil,
			},
			wantErr: false,
		},
		{
			name: "When VideoDTO is right",
			args: args{
				crud.VideoDTO{
					Title: fakeDoesNotExistTitle,
				},
			},
			want: returns{
				models.Video{
					Title: fakeDoesNotExistTitle,
				},
				nil,
			},
			wantErr: false,
		},
	}
	db, err := sql.Open(cfg.DBDrive, cfg.DBConnStr)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("test: failed to close DB: %v\n", err)
		}
	}()
	ctx := context.Background()
	err = fakeExistVideo.InsertG(ctx, boil.Infer())
	if err != nil {
		t.Errorf("test: insert video: %s", err)
		return
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.AddVideo(tt.args.videoDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddVideo() error: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(err, tt.want.err) {
				t.Errorf("AddVideo() got: %v, want: %v", err, tt.want.err)
				return
			}
		})
	}
}

func TestRepository_GetVideos(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeVideos)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	maximum := len(testdata.FakeVideos)
	type args struct {
		limit int
	}
	type returns struct {
		videos models.VideoSlice
		e      error
		amount int
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name:    "When limit is negative",
			args:    args{-1},
			want:    returns{nil, nil, 0},
			wantErr: false,
		},
		{
			name:    "When limit is zero",
			args:    args{0},
			want:    returns{nil, nil, 0},
			wantErr: false,
		},
		{
			name:    "When limit is less then the maximum",
			args:    args{maximum - 1},
			want:    returns{nil, nil, maximum - 1},
			wantErr: false,
		},
		{
			name:    "When limit is equal the maximum",
			args:    args{maximum},
			want:    returns{nil, nil, maximum},
			wantErr: false,
		},
		{
			name:    "When limit is more then the maximum",
			args:    args{maximum + 1},
			want:    returns{nil, nil, maximum},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repository.GetVideos(tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetVideos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want.amount {
				t.Errorf("GetVideos() len(got): %v, want: %d", len(got), tt.want.amount)
			}
		})
	}
}

func TestRepository_FetchVideo(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeVideos)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	fakeExistTitle := testdata.FakeVideos[0].Title
	const (
		fakeDoesNotExistTitle = "fakeDoesNotExistTitle"
	)
	type args struct {
		title string
	}
	type returns struct {
		video models.Video
		e     error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "When title is not found",
			args: args{fakeDoesNotExistTitle},
			want: returns{
				models.Video{},
				sql.ErrNoRows,
			},
			wantErr: true,
		},
		{
			name: "When title is found",
			args: args{fakeExistTitle},
			want: returns{
				models.Video{
					Title: fakeExistTitle,
				},
				nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repository.FetchVideo(tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchVideo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Title != tt.want.video.Title {
				t.Errorf("FetchVideo() got: %q, want: %q", got.Title, tt.want.video.Title)
			}
		})
	}
}

func TestRepository_RemoveVideo(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeVideos)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	fakeExistTitle := testdata.FakeVideos[0].Title
	const (
		fakeDoesNotExistTitle = "fakeDoesNotExistTitle"
	)
	type fields struct {
		ctx context.Context
	}
	type args struct {
		title string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    error
		wantErr bool
	}{
		{
			name:    "When title is not found",
			args:    args{fakeDoesNotExistTitle},
			want:    sql.ErrNoRows,
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
			err := repository.RemoveVideo(tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveVideo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != tt.want {
				t.Errorf("RemoveVideo() got: %s, want: %q", err, tt.want)
			}
		})
	}
}

func TestRepository_UpdateVideo(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeVideos)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	fakeExistTitle := testdata.FakeVideos[0].Title
	const (
		fakeDoesNotExistTitle     = "fakeDoesNotExistTitle"
		fakeNewDoestNotExistTitle = "new_action"
	)
	type fields struct {
		ctx context.Context
	}
	type args struct {
		title    string
		videoDTO crud.VideoDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "When title to update doesn't exist",
			args: args{
				fakeDoesNotExistTitle,
				crud.VideoDTO{
					Title: fakeNewDoestNotExistTitle,
				},
			},
			want:    sql.ErrNoRows,
			wantErr: true,
		},
		{
			name: "When title exists and VideoDTO is right",
			args: args{
				fakeExistTitle,
				crud.VideoDTO{
					Title: fakeNewDoestNotExistTitle,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.UpdateVideo(tt.args.title, tt.args.videoDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateVideo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("UpdateVideo() got: %v, want: %v", err, tt.want)
			}
		})
	}
}

func TestVideo_isValidUUIDHook(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(nil)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	type args struct {
		video models.Video
	}
	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "When UUID is not validated",
			args: args{
				models.Video{
					ID:    "fakeUUIDIsNotValidated",
					Title: faker.Name(),
				},
			},
			want:    fmt.Errorf("%s %w", "UUID", logger.ErrIsNotValidated),
			wantErr: true,
		},
		{
			name: "When UUID is validated",
			args: args{
				models.Video{
					ID:           uuid.New().String(),
					Title:        faker.Name(),
					Description:  faker.Sentence(),
					YearLaunched: 2020,
					Opened:       null.Bool{Bool: true, Valid: true},
					Rating:       int16(crud.TwelveRating),
					Duration:     150,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.video.InsertG(repository.ctx, boil.Infer())
			if (err != nil) != tt.wantErr {
				t.Errorf("isValidUUIDHook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("UpdateCategory() got: %v, want: %v", err, tt.want)
			}
		})
	}
}