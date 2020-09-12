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

var (
	fakeCtx  = context.Background()
	fakeDesc = "fakeDesc"
)

func TestRepository_CreateCategory(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(nil)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	const (
		fakeExistCategoryName        = "fakeExistCategoryName"
		fakeExistGenreName           = "fakeExistGenreName"
		fakeDoesNotExistCategoryName = "fakeDoesNotExistCategoryName"
		fakeDoesNotExistGenreName    = "fakeDoesNotExistGenreName"
	)
	fakeExistCategory := service.Category{Name: fakeExistCategoryName}
	fakeDoesNotExistGenre := service.Genre{Name: fakeDoesNotExistGenreName}
	fakeExistGenre := service.Genre{
		Name: fakeExistGenreName,
	}
	fakeExistGenresOfCategory := []service.GenreOfCategory{
		{
			Name: fakeExistGenre.Name,
		},
	}
	fakeDoesNotExistGenresOfCategory := []service.GenreOfCategory{
		{Name: fakeDoesNotExistGenre.Name},
	}
	if err := repository.CreateCategory(fakeCtx, fakeExistCategory); err != nil {
		t.Errorf("test: insert category: %s", err)
		return
	}
	if err := repository.CreateGenre(fakeCtx, fakeExistGenre); err != nil {
		t.Errorf("test: insert genre: %s", err)
		return
	}
	fetchGenre, err := repository.FetchGenre(fakeCtx, fakeExistGenre.Name)
	if err != nil {
		t.Errorf("test: insert category: %s", err)
		return
	}
	fakeExistGenre.Id = fetchGenre.Id
	fakeExistGenresOfCategory[0].Id = fetchGenre.Id
	type args struct {
		ctx      context.Context
		category service.Category
	}
	type returns struct {
		category models.Category
		err      error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "When name in Category already exists",
			args: args{
				fakeCtx,
				service.Category{
					Name: fakeExistCategoryName,
				},
			},
			want:    returns{models.Category{}, fmt.Errorf("name '%s' %w", fakeExistCategoryName, logger.ErrAlreadyExists)},
			wantErr: true,
		},
		{
			name: "When VideoDTO is with wrong genres",
			args: args{
				fakeCtx,
				service.Category{
					Name:        fakeDoesNotExistCategoryName,
					Description: &fakeDesc,
					Genres:      &fakeDoesNotExistGenresOfCategory,
				},
			},
			want:    returns{err: fmt.Errorf("none genre is %w", logger.ErrNotFound)},
			wantErr: true,
		},
		{
			name: "When Category is right",
			args: args{
				fakeCtx,
				service.Category{
					Name:   fakeDoesNotExistCategoryName,
					Genres: &fakeExistGenresOfCategory,
				},
			},
			want: returns{
				models.Category{
					Name: fakeDoesNotExistCategoryName,
				},
				nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.CreateCategory(tt.args.ctx, tt.args.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCategory() error: %v, wantErr %v", err, tt.wantErr)
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

func TestRepository_GetCategories(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeCategories)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	maximum := len(testdata.FakeCategories)
	type args struct {
		ctx   context.Context
		limit int
	}
	type returns struct {
		categories models.CategorySlice
		e          error
		amount     int
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
			got, err := repository.GetCategories(tt.args.ctx, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want.amount {
				t.Errorf("GetCategories() len(got): %v, want: %d", len(got), tt.want.amount)
			}
		})
	}
}

func TestRepository_FetchCategory(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeCategories)
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
		category models.Category
		e        error
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
				models.Category{},
				sql.ErrNoRows,
			},
			wantErr: true,
		},
		{
			name: "When name is found",
			args: args{fakeCtx, fakeExistName},
			want: returns{
				models.Category{
					Name: fakeExistName,
				},
				nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repository.FetchCategory(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Name != tt.want.category.Name {
				t.Errorf("FetchCategory() got: %q, want: %q", got.Name, tt.want.category.Name)
			}
		})
	}
}

func TestRepository_RemoveCategory(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeCategories)
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
			err := repository.RemoveCategory(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != tt.want {
				t.Errorf("RemoveCategory() got: %s, want: %q", err, tt.want)
			}
		})
	}
}

func TestRepository_UpdateCategory(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(nil)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	const (
		fakeExistCategoryName            = "fakeExistCategoryName"
		fakeNewExistCategoryName         = "fakeNewExistCategoryName"
		fakeDoesNotExistCategoryName     = "fakeDoesNotExistCategoryName"
		fakeDoesNotExistGenreName        = "fakeDoesNotExistGenreName"
		fakeNewDoestNotExistCategoryName = "fakeNewDoestNotExistCategoryName"
		fakeExistGenreName               = "fakeExistGenreName"
	)
	fakeExitUUID := uuid.New().String()
	fakeDoesNotExitUUID := uuid.New().String()
	fakeDoesNotExistGenre := service.Genre{Id: &fakeDoesNotExitUUID, Name: fakeDoesNotExistGenreName}
	fakeExistCategory := service.Category{Name: fakeExistCategoryName}
	fakeNewExistCategory := service.Category{Name: fakeNewExistCategoryName}
	fakeExistGenre := service.Genre{Id: &fakeExitUUID, Name: fakeExistGenreName}
	fakeExistGenres := []service.GenreOfCategory{
		{
			Id:   fakeExistGenre.Id,
			Name: fakeExistGenre.Name,
		},
	}
	fakeDoesNotExistGenres := []service.GenreOfCategory{
		{
			Id:   fakeDoesNotExistGenre.Id,
			Name: fakeDoesNotExistGenre.Name,
		},
	}
	if err := repository.CreateCategory(fakeCtx, fakeExistCategory); err != nil {
		t.Errorf("test: insert category: %s", err)
		return
	}
	if err := repository.CreateCategory(fakeCtx, fakeNewExistCategory); err != nil {
		t.Errorf("test: insert category: %s", err)
		return
	}
	if err := repository.CreateGenre(fakeCtx, fakeExistGenre); err != nil {
		t.Errorf("test: insert genre: %s", err)
		return
	}
	type args struct {
		ctx         context.Context
		name        string
		categoryDTO service.Category
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
				fakeDoesNotExistCategoryName,
				service.Category{
					Name:   fakeNewDoestNotExistCategoryName,
					Genres: &fakeExistGenres,
				},
			},
			want:    fmt.Errorf("%s: %w", fakeDoesNotExistCategoryName, logger.ErrNotFound),
			wantErr: true,
		},
		{
			name: "When name in Category already exists",
			args: args{
				fakeCtx,
				fakeExistCategoryName,
				service.Category{
					Name:   fakeNewExistCategoryName,
					Genres: &fakeExistGenres,
				},
			},
			want:    fmt.Errorf("%s %w", fakeNewExistCategoryName, logger.ErrAlreadyExists),
			wantErr: true,
		},
		{
			name: "When Category is with wrong genres",
			args: args{
				fakeCtx,
				fakeExistCategoryName,
				service.Category{
					Name:        fakeDoesNotExistCategoryName,
					Description: &fakeDesc,
					Genres:      &fakeDoesNotExistGenres,
				},
			},
			want:    fmt.Errorf("none genre is %w", logger.ErrNotFound),
			wantErr: true,
		},
		{
			name: "When everything is right",
			args: args{
				fakeCtx,
				fakeExistCategoryName,
				service.Category{
					Name:        fakeNewDoestNotExistCategoryName,
					Description: &fakeDesc,
					Genres:      &fakeExistGenres,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.UpdateCategory(tt.args.ctx, tt.args.name, tt.args.categoryDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCategory() got: %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && err.Error() != tt.want.Error() {
				t.Errorf("UpdateCategory() got: %v, want: %v", err, tt.want)
			}
		})
	}
}

func TestCategory_isValidUUIDHook(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeCategories)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	type args struct {
		category models.Category
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
				models.Category{
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
				models.Category{
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
			err := tt.args.category.InsertG(repository.ctx, boil.Infer())
			if (err != nil) != tt.wantErr {
				t.Errorf("isValidUUIDCategoryHook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("UpdateCategory() got: %v, want: %v", err, tt.want)
			}
		})
	}
}
