// +build integration

package sqlboiler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud/service"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/testdata"
)

var (
	fakeYearLaunched = new(int16)
	fakeDuration     = new(int16)
	fakeRating       = new(service.VideoRating)
)

func TestRepository_AddVideo(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(nil)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	const (
		fakeDoesNotExistTitle = "fakeDoesNotExistTitle"
		fakeOpened            = false
		fakeCategoryIndex     = 0
		fakeGenreIndex        = 0
	)
	*fakeYearLaunched = 2020
	*fakeDuration = 90
	*fakeRating = service.TwelveRating
	fakeExistCategoryDTO := testdata.FakeCategoriesDTO[fakeCategoryIndex]
	fakeDoesNotExistCategoryDTO := service.Category{Name: faker.FirstName(), Description: &fakeDesc}
	fakeExistGenreDTO := testdata.FakeGenresDTO[fakeGenreIndex]
	fakeDoesNotExistGenreDTO := service.Genre{Name: faker.FirstName()}
	if err := repository.CreateCategory(fakeCtx, fakeExistCategoryDTO); err != nil {
		t.Errorf("test: insert category: %s", err)
		return
	}
	if err := repository.CreateGenre(fakeCtx, fakeExistGenreDTO); err != nil {
		t.Errorf("test: insert genre: %s", err)
		return
	}
	type args struct {
		videoDTO service.VideoDTO
	}
	type returns struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "When File in VideoDTO is omitted",
			args: args{
				service.VideoDTO{
					Title:        fakeDoesNotExistTitle,
					Description:  fakeDesc,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
					Genres:       []service.Genre{fakeDoesNotExistGenreDTO},
					Categories:   []service.Category{fakeDoesNotExistCategoryDTO},
				},
			},
			want:    returns{logger.ErrNotFound},
			wantErr: true,
		},
		{
			name: "When VideoDTO is with wrong categories and genres",
			args: args{
				service.VideoDTO{
					Title:        fakeDoesNotExistTitle,
					Description:  fakeDesc,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
					Genres:       []service.Genre{fakeDoesNotExistGenreDTO},
					Categories:   []service.Category{fakeDoesNotExistCategoryDTO},
				},
			},
			want:    returns{logger.ErrNotFound},
			wantErr: true,
		},
		{
			name: "When VideoDTO is right",
			args: args{service.VideoDTO{
				Title:        fakeDoesNotExistTitle,
				Description:  fakeDesc,
				YearLaunched: fakeYearLaunched,
				Opened:       fakeOpened,
				Rating:       fakeRating,
				Duration:     fakeDuration,
				Genres:       []service.Genre{fakeExistGenreDTO},
				Categories:   []service.Category{fakeExistCategoryDTO},
			}},
			want:    returns{nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repository.AddVideo(tt.args.videoDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddVideo() got: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !errors.Is(err, tt.want.err) {
					t.Errorf("AddVideo() got: %v, want: %v", err, tt.want.err)
				}
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
	maximum := testdata.FakeVideosLength
	fakeVideosSlice := testdata.FakeVideoSlice
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
			want:    returns{fakeVideosSlice[:maximum-1], nil, maximum - 1},
			wantErr: false,
		},
		{
			name:    "When limit is equal the maximum",
			args:    args{maximum},
			want:    returns{fakeVideosSlice, nil, maximum},
			wantErr: false,
		},
		{
			name:    "When limit is more then the maximum",
			args:    args{maximum + 1},
			want:    returns{fakeVideosSlice, nil, maximum},
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
			gotRemovedTimes := removeUnwantedStuff(got)
			videosRemovedTimes := removeUnwantedStuff(tt.want.videos)
			assert.Equal(t, gotRemovedTimes, videosRemovedTimes, "they should be equal")
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
	fakeExistVideo := testdata.FakeVideos[0]
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
			args: args{fakeExistVideo.Title},
			want: returns{
				fakeExistVideo,
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
			gotRemovedTimes := removeUnwantedStuff(got)
			videoRemovedTimes := removeUnwantedStuff(tt.want.video)
			assert.Equal(t, gotRemovedTimes, videoRemovedTimes, "they should be equal")
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
	const (
		fakeDoesNotExistTitle     = "fakeDoesNotExistTitle"
		fakeExistGenreName        = "fakeExistGenreName"
		fakeExistCategoryName     = "fakeExistCategoryName"
		fakeOpened                = false
		fakeNewDoestNotExistTitle = "fakeNewDoestNotExistTitle"
	)
	*fakeYearLaunched = 2020
	*fakeDuration = 90
	*fakeRating = service.TwelveRating
	fakeDoesNotExistCategoryDTO := service.Category{Name: faker.FirstName(), Description: &fakeDesc}
	fakeDoesNotExistGenreDTO := service.Genre{Name: faker.FirstName()}
	fakeExistCategoryDTO := service.Category{Name: fakeExistCategoryName}
	fakeExistGenreDTO := service.Genre{Name: fakeExistGenreName}
	fakeExistTitle := strings.ToLower(testdata.FakeVideos[0].Title)
	fakeExistVideoDTO := service.VideoDTO{
		Title:        fakeExistTitle,
		Description:  fakeDesc,
		YearLaunched: fakeYearLaunched,
		Opened:       fakeOpened,
		Rating:       fakeRating,
		Duration:     fakeDuration,
		Genres:       []service.Genre{fakeExistGenreDTO},
		Categories:   []service.Category{fakeExistCategoryDTO},
	}
	if err := repository.CreateCategory(fakeCtx, fakeExistCategoryDTO); err != nil {
		t.Errorf("test: insert category: %s", err)
		return
	}
	if err := repository.CreateGenre(fakeCtx, fakeExistGenreDTO); err != nil {
		t.Errorf("test: insert genre: %s", err)
		return
	}
	if _, err := repository.AddVideo(fakeExistVideoDTO); err != nil {
		t.Errorf("test: insert category: %s", err)
		return
	}
	type fields struct {
		ctx context.Context
	}
	type args struct {
		title    string
		videoDTO service.VideoDTO
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
				service.VideoDTO{
					Title: fakeNewDoestNotExistTitle,
				},
			},
			want:    sql.ErrNoRows,
			wantErr: true,
		},
		{
			name: "When VideoDTO is with wrong genres",
			args: args{
				fakeExistTitle,
				service.VideoDTO{
					Title:        fakeDoesNotExistTitle,
					Description:  fakeDesc,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
					Genres:       []service.Genre{fakeDoesNotExistGenreDTO},
					Categories:   []service.Category{fakeDoesNotExistCategoryDTO},
				},
			},
			want:    fmt.Errorf("none category is %w", logger.ErrNotFound),
			wantErr: true,
		},
		{
			name: "When everything is right",
			args: args{
				fakeExistTitle,
				service.VideoDTO{
					Title:        fakeDoesNotExistTitle,
					Description:  fakeDesc,
					YearLaunched: fakeYearLaunched,
					Opened:       fakeOpened,
					Rating:       fakeRating,
					Duration:     fakeDuration,
					Genres:       []service.Genre{fakeExistGenreDTO},
					Categories:   []service.Category{fakeExistCategoryDTO},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := repository.UpdateVideo(tt.args.title, tt.args.videoDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateVideo() got: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.want.Error() {
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
			name: "When ID is not validated",
			args: args{
				models.Video{
					ID:    "fakeUUIDIsNotValidated",
					Title: faker.Name(),
				},
			},
			want:    fmt.Errorf("%s %w", "ID", logger.ErrIsNotValidated),
			wantErr: true,
		},
		{
			name: "When ID is validated",
			args: args{
				models.Video{
					ID:           uuid.New().String(),
					Title:        faker.Name(),
					Description:  faker.Sentence(),
					YearLaunched: 2020,
					Opened:       null.Bool{Bool: true, Valid: true},
					Rating:       int16(service.TwelveRating),
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

func removeUnwantedStuff(i interface{}) interface{} {
	if i == nil || reflect.ValueOf(i).IsZero() {
		return i
	}
	funcRemoveTimes := func(i interface{}) {
		value := reflect.ValueOf(i)
		for _, fieldName := range [3]string{"CreatedAt", "DeletedAt", "UpdatedAt"} {
			field := reflect.Indirect(value).FieldByName(fieldName)
			if field.IsValid() && !field.IsZero() {
				field.Set(
					reflect.ValueOf(null.Time{
						Time:  time.Time{},
						Valid: false,
					}),
				)
			}
		}
	}
	switch v := i.(type) {
	case models.VideoSlice:
		for _, video := range v {
			funcRemoveTimes(video)
			if video.R == nil {
				continue
			}
			for _, g := range video.R.Genres {
				funcRemoveTimes(g)
				g.R = nil
			}
			for _, c := range video.R.Categories {
				funcRemoveTimes(c)
				c.R = nil
			}
		}
		return v
	case models.Video:
		funcRemoveTimes(&v)
		if v.R == nil {
			return v
		}
		for _, c := range v.R.Categories {
			funcRemoveTimes(c)
			c.R = nil
		}
		for _, g := range v.R.Genres {
			funcRemoveTimes(g)
			g.R = nil
		}
		return v
	}
	return nil
}
