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
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud/service"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/testdata"
)

func TestRepository_CreateGenre(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(nil)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	const (
		fakeExistGenreName           = "fakeExistGenreName"
		fakeExistCategoryName        = "fakeExistCategoryName"
		fakeDoesNotExistGenreName    = "fakeDoesNotExistGenreName"
		fakeDoesNotExistCategoryName = "fakeDoesNotExistCategoryName"
	)
	fakeExistGenre := service.Genre{Name: fakeExistGenreName}
	fakeDoesNotExistCategory := service.Category{Name: fakeDoesNotExistCategoryName}
	fakeExistCategory := service.Category{Name: fakeExistCategoryName}
	fakeExistCategories := []service.CategoryOfGenre{
		{
			Id:   fakeExistCategory.Id,
			Name: fakeExistCategory.Name,
		},
	}
	fakeDoesNotExistCategories := []service.CategoryOfGenre{
		{
			Id:   fakeDoesNotExistCategory.Id,
			Name: fakeDoesNotExistCategory.Name,
		},
	}
	if err := repository.CreateGenre(fakeCtx, fakeExistGenre); err != nil {
		t.Errorf("test: insert genre: %s", err)
		return
	}
	if err := repository.CreateCategory(fakeCtx, fakeExistCategory); err != nil {
		t.Errorf("test: insert genre: %s", err)
		return
	}
	type args struct {
		ctx      context.Context
		genreDTO service.Genre
	}
	type returns struct {
		genre models.Genre
		err   error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "When name in Genre already exists",
			args: args{
				fakeCtx,
				service.Genre{
					Name: fakeExistGenreName,
				}},
			want:    returns{models.Genre{}, fmt.Errorf("name '%s' %w", fakeExistGenreName, logger.ErrAlreadyExists)},
			wantErr: true,
		},
		{
			name: "When VideoDTO is with wrong genres",
			args: args{
				fakeCtx,
				service.Genre{
					Name:       fakeDoesNotExistGenreName,
					Categories: &fakeDoesNotExistCategories,
				},
			},
			want:    returns{err: fmt.Errorf("none category is %w", logger.ErrNotFound)},
			wantErr: true,
		},
		{
			name: "When Genre is right",
			args: args{
				fakeCtx,
				service.Genre{
					Name:       fakeDoesNotExistGenreName,
					Categories: &fakeExistCategories,
				},
			},
			want: returns{
				models.Genre{
					Name: fakeDoesNotExistGenreName,
				},
				nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.CreateGenre(tt.args.ctx, tt.args.genreDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateGenre() error: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if err != nil && err.Error() != tt.want.err.Error() {
					t.Errorf("AddVideo() got: %v, want: %v", err, tt.want.err)
				}
				return
			}
		})
	}
}

func TestRepository_GetGenres(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeGenres)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	maximum := len(testdata.FakeGenres)
	type args struct {
		ctx   context.Context
		limit int
	}
	type returns struct {
		genres models.GenreSlice
		err    error
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
			args:    args{fakeCtx, -1},
			want:    returns{nil, nil, 0},
			wantErr: false,
		},
		{
			name:    "When limit is zero",
			args:    args{fakeCtx, 0},
			want:    returns{nil, nil, 0},
			wantErr: false,
		},
		{
			name:    "When limit is less then the maximum",
			args:    args{fakeCtx, maximum - 1},
			want:    returns{nil, nil, maximum - 1},
			wantErr: false,
		},
		{
			name:    "When limit is equal the maximum",
			args:    args{fakeCtx, maximum},
			want:    returns{nil, nil, maximum},
			wantErr: false,
		},
		{
			name:    "When limit is more then the maximum",
			args:    args{fakeCtx, maximum + 1},
			want:    returns{nil, nil, maximum},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repository.GetGenres(tt.args.ctx, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGenres() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want.amount {
				t.Errorf("GetGenres() len(got): %v, want: %d", len(got), tt.want.amount)
			}
		})
	}
}

func TestRepository_FetchGenre(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeGenres)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	const (
		fakeExistName        = "action"
		fakeDoesNotExistName = "fakeDoesNotExistName"
	)
	type args struct {
		ctx  context.Context
		name string
	}
	type returns struct {
		genre models.Genre
		e     error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "When name is not found",
			args: args{fakeCtx, fakeDoesNotExistName},
			want: returns{
				models.Genre{},
				sql.ErrNoRows,
			},
			wantErr: true,
		},
		{
			name: "When name is found",
			args: args{fakeCtx, fakeExistName},
			want: returns{
				models.Genre{
					Name: fakeExistName,
				},
				nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repository.FetchGenre(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchGenre() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got.Name != tt.want.genre.Name {
				t.Errorf("FetchGenre() got: %q, want: %q", got.Name, tt.want.genre.Name)
			}
		})
	}
}

func TestRepository_RemoveGenre(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeGenres)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	const (
		fakeExistName        = "action"
		fakeDoesNotExistName = "fakeDoesNotExistName"
	)
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name:    "When name is not found",
			args:    args{fakeCtx, fakeDoesNotExistName},
			want:    sql.ErrNoRows,
			wantErr: true,
		},
		{
			name:    "When name is found",
			args:    args{fakeCtx, fakeExistName},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.RemoveGenreByName(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveGenreByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != tt.want {
				t.Errorf("RemoveGenreByName() got: %s, want: %q", err, tt.want)
			}
		})
	}
}

func TestRepository_UpdateGenre(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(nil)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	const (
		fakeExistGenreName            = "fakeExistGenreName"
		fakeNewExistGenreName         = "fakeNewExistGenreName"
		fakeDoesNotExistGenreName     = "fakeDoesNotExistGenreName"
		fakeDoesNotExistCategoryName  = "fakeDoesNotExistCategoryName"
		fakeNewDoestNotExistGenreName = "fakeNewDoestNotExistGenreName"
		fakeExistCategoryName         = "fakeExistCategoryName"
	)
	fakeDoesNotExistCategory := service.Category{Name: fakeDoesNotExistCategoryName}
	fakeExistGenre := service.Genre{Name: fakeExistGenreName}
	fakeNewExistGenre := service.Genre{Name: fakeNewExistGenreName}
	fakeExistCategory := service.Category{Name: fakeExistCategoryName}
	fakeExistCategoriesOfGenre := []service.CategoryOfGenre{
		{
			Id:          fakeExistCategory.Id,
			Name:        fakeExistCategory.Name,
			Description: fakeExistCategory.Description,
		},
	}
	fakeDoesNotExistCategoriesOfGenre := []service.CategoryOfGenre{
		{
			Id:          fakeDoesNotExistCategory.Id,
			Name:        fakeDoesNotExistCategory.Name,
			Description: fakeDoesNotExistCategory.Description,
		},
	}

	if err := repository.CreateGenre(fakeCtx, fakeExistGenre); err != nil {
		t.Errorf("test: insert genre: %s", err)
		return
	}
	if err := repository.CreateGenre(fakeCtx, fakeNewExistGenre); err != nil {
		t.Errorf("test: insert genre: %s", err)
		return
	}
	if err := repository.CreateCategory(fakeCtx, fakeExistCategory); err != nil {
		t.Errorf("test: insert category: %s", err)
		return
	}
	type args struct {
		ctx      context.Context
		name     string
		genreDTO service.Genre
	}
	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "When name to update doesn't exist",
			args: args{
				fakeCtx,
				fakeDoesNotExistGenreName,
				service.Genre{
					Name:       fakeNewDoestNotExistGenreName,
					Categories: &fakeExistCategoriesOfGenre,
				},
			},
			want:    sql.ErrNoRows,
			wantErr: true,
		},
		{
			name: "When name in Genre already exists",
			args: args{
				fakeCtx,
				fakeExistGenreName,
				service.Genre{
					Name:       fakeNewExistGenreName,
					Categories: &fakeExistCategoriesOfGenre,
				},
			},
			want:    fmt.Errorf("%s %w", fakeNewExistGenreName, logger.ErrAlreadyExists),
			wantErr: true,
		},
		{
			name: "When Genre is with wrong genres",
			args: args{
				fakeCtx,
				fakeExistGenreName,
				service.Genre{
					Name:       fakeDoesNotExistGenreName,
					Categories: &fakeDoesNotExistCategoriesOfGenre,
				},
			},
			want:    fmt.Errorf("none category is %w", logger.ErrNotFound),
			wantErr: true,
		},
		{
			name: "When everything is right",
			args: args{
				fakeCtx,
				fakeExistGenreName,
				service.Genre{
					Name:       fakeNewDoestNotExistGenreName,
					Categories: &fakeExistCategoriesOfGenre,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.UpdateGenre(tt.args.ctx, tt.args.name, tt.args.genreDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateGenre() got: %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && err.Error() != tt.want.Error() {
				t.Errorf("UpdateGenre() got: %v, want: %v", err, tt.want)
			}
		})
	}
}

func TestGenre_isValidUUIDHook(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeGenres)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	type args struct {
		genre models.Genre
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
				models.Genre{
					ID:   "fakeUUIDIsNotValidated",
					Name: faker.FirstName(),
				},
			},
			want:    fmt.Errorf("%s %w", "ID", logger.ErrIsNotValidated),
			wantErr: true,
		},
		{
			name: "When ID is validated",
			args: args{
				models.Genre{
					ID:   uuid.New().String(),
					Name: faker.FirstName(),
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.genre.InsertG(repository.ctx, boil.Infer())
			if (err != nil) != tt.wantErr {
				t.Errorf("isValidUUIDGenreHook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("UpdateGenre() got: %v, want: %v", err, tt.want)
			}
		})
	}
}
