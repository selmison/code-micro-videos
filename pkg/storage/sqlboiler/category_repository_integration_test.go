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
	"github.com/selmison/code-micro-videos/pkg/crud"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/testdata"
)

func TestRepository_AddCategory(t *testing.T) {
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
		fakeDesc                     = "fakeDesc"
	)
	fakeExistCategoryDTO := crud.CategoryDTO{Name: fakeExistCategoryName}
	fakeDoesNotExistGenreDTO := crud.GenreDTO{Name: fakeDoesNotExistGenreName}
	fakeExistGenreDTO := crud.GenreDTO{Name: fakeExistGenreName}
	if err := repository.AddCategory(fakeExistCategoryDTO); err != nil {
		t.Errorf("test: insert category: %s", err)
		return
	}
	if err := repository.AddGenre(fakeExistGenreDTO); err != nil {
		t.Errorf("test: insert genre: %s", err)
		return
	}
	type args struct {
		categoryDTO crud.CategoryDTO
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
			name: "When name in CategoryDTO already exists",
			args: args{crud.CategoryDTO{
				Name: fakeExistCategoryName,
			}},
			want:    returns{models.Category{}, fmt.Errorf("name '%s' %w", fakeExistCategoryName, logger.ErrAlreadyExists)},
			wantErr: true,
		},
		{
			name: "When VideoDTO is with wrong genres",
			args: args{
				crud.CategoryDTO{
					Name:        fakeDoesNotExistCategoryName,
					Description: fakeDesc,
					Genres:      []crud.GenreDTO{fakeDoesNotExistGenreDTO},
				},
			},
			want:    returns{err: fmt.Errorf("none genre is %w", logger.ErrNotFound)},
			wantErr: true,
		},
		{
			name: "When CategoryDTO is right",
			args: args{
				crud.CategoryDTO{
					Name:   fakeDoesNotExistCategoryName,
					Genres: []crud.GenreDTO{fakeExistGenreDTO},
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
			err := repository.AddCategory(tt.args.categoryDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCategory() error: %v, wantErr %v", err, tt.wantErr)
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
			got, err := repository.GetCategories(tt.args.limit)
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
			args: args{fakeDoesNotExistName},
			want: returns{
				models.Category{},
				sql.ErrNoRows,
			},
			wantErr: true,
		},
		{
			name: "When name is found",
			args: args{fakeExistName},
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
			got, err := repository.FetchCategory(tt.args.name)
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
	type fields struct {
		ctx context.Context
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    error
		wantErr bool
	}{
		{
			name:    "When name is not found",
			args:    args{fakeDoesNotExistName},
			want:    sql.ErrNoRows,
			wantErr: true,
		},
		{
			name:    "When name is found",
			args:    args{fakeExistName},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.RemoveCategory(tt.args.name)
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
		fakeDesc                         = "fakeDesc"
	)
	fakeDoesNotExistGenreDTO := crud.GenreDTO{Name: fakeDoesNotExistGenreName}
	fakeExistCategoryDTO := crud.CategoryDTO{Name: fakeExistCategoryName}
	fakeNewExistCategoryDTO := crud.CategoryDTO{Name: fakeNewExistCategoryName}
	fakeExistGenreDTO := crud.GenreDTO{Name: fakeExistGenreName}
	if err := repository.AddCategory(fakeExistCategoryDTO); err != nil {
		t.Errorf("test: insert category: %s", err)
		return
	}
	if err := repository.AddCategory(fakeNewExistCategoryDTO); err != nil {
		t.Errorf("test: insert category: %s", err)
		return
	}
	if err := repository.AddGenre(fakeExistGenreDTO); err != nil {
		t.Errorf("test: insert genre: %s", err)
		return
	}
	type args struct {
		name        string
		categoryDTO crud.CategoryDTO
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
				fakeDoesNotExistCategoryName,
				crud.CategoryDTO{
					Name:   fakeNewDoestNotExistCategoryName,
					Genres: []crud.GenreDTO{fakeExistGenreDTO},
				},
			},
			want:    sql.ErrNoRows,
			wantErr: true,
		},
		{
			name: "When name in CategoryDTO already exists",
			args: args{
				fakeExistCategoryName,
				crud.CategoryDTO{
					Name:   fakeNewExistCategoryName,
					Genres: []crud.GenreDTO{fakeExistGenreDTO},
				},
			},
			want:    fmt.Errorf("%s %w", fakeNewExistCategoryName, logger.ErrAlreadyExists),
			wantErr: true,
		},
		{
			name: "When CategoryDTO is with wrong genres",
			args: args{
				fakeExistCategoryName,
				crud.CategoryDTO{
					Name:        fakeDoesNotExistCategoryName,
					Description: fakeDesc,
					Genres:      []crud.GenreDTO{fakeDoesNotExistGenreDTO},
				},
			},
			want:    fmt.Errorf("none genre is %w", logger.ErrNotFound),
			wantErr: true,
		},
		{
			name: "When everything is right",
			args: args{
				fakeExistCategoryName,
				crud.CategoryDTO{
					Name:        fakeNewDoestNotExistCategoryName,
					Description: fakeDesc,
					Genres:      []crud.GenreDTO{fakeExistGenreDTO},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.UpdateCategory(tt.args.name, tt.args.categoryDTO)
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
			name: "When UUID is not validated",
			args: args{
				models.Category{
					ID:   "fakeUUIDIsNotValidated",
					Name: faker.FirstName(),
				},
			},
			want:    fmt.Errorf("%s %w", "UUID", logger.ErrIsNotValidated),
			wantErr: true,
		},
		{
			name: "When UUID is validated",
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
