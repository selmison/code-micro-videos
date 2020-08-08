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

	"github.com/selmison/code-micro-videos/config"
	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/testdata"
)

func TestRepository_AddCategory(t *testing.T) {
	const (
		fakeExistName        = "action"
		fakeDoesNotExistName = "fakeDoesNotExistName"
	)
	fakeExistCategory := models.Category{
		ID:   uuid.New().String(),
		Name: fakeExistName,
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
				Name: fakeExistName,
			}},
			want:    returns{models.Category{}, fmt.Errorf("name '%s' %w", fakeExistName, logger.ErrAlreadyExists)},
			wantErr: true,
		},
		{
			name: "When CategoryDTO is right",
			args: args{
				crud.CategoryDTO{
					Name: fakeDoesNotExistName,
				},
			},
			want: returns{
				models.Category{
					Name: fakeDoesNotExistName,
				},
				nil,
			},
			wantErr: false,
		},
	}
	db, err := sql.Open(config.DBDrive, dbConnStr)
	if err != nil {
		t.Errorf("test: failed to open DB: %v\n", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("test: failed to close DB: %v\n", err)
		}
	}()
	ctx := context.Background()
	r := NewRepository(ctx, db)
	err = fakeExistCategory.InsertG(ctx, boil.Infer())
	if err != nil {
		t.Errorf("test: insert category: %s", err)
		return
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.AddCategory(tt.args.categoryDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCategory() error: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(err, tt.want.err) {
				t.Errorf("AddCategory() got: %v, want: %v", err, tt.want.err)
				return
			}
		})
	}
	if err := config.ClearCategoriesTable(config.DBDrive, dbConnStr); err != nil {
		t.Errorf("test: clear categories table: %v", err)
	}
}

func TestRepository_GetCategories(t *testing.T) {
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
	db, err := sql.Open(config.DBDrive, dbConnStr)
	if err != nil {
		t.Errorf("test: failed to open DB: %v\n", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("test: failed to close DB: %v\n", err)
		}
	}()
	ctx := context.Background()
	r := NewRepository(ctx, db)
	for _, g := range testdata.FakeCategories {
		err = g.InsertG(ctx, boil.Infer())
		if err != nil {
			t.Errorf("test: insert category: %s", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.GetCategories(tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want.amount {
				t.Errorf("GetCategories() len(got): %v, want: %d", len(got), tt.want.amount)
			}
		})
	}
	if err := config.ClearCategoriesTable(config.DBDrive, dbConnStr); err != nil {
		t.Errorf("test: clear categories table: %v", err)
	}
}

func TestRepository_FetchCategory(t *testing.T) {
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
	db, err := sql.Open(config.DBDrive, dbConnStr)
	if err != nil {
		t.Errorf("test: failed to open DB: %v\n", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("test: failed to close DB: %v\n", err)
		}
	}()
	ctx := context.Background()
	r := NewRepository(ctx, db)
	for _, g := range testdata.FakeCategories {
		err = g.InsertG(ctx, boil.Infer())
		if err != nil {
			t.Errorf("test: insert category: %s", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.FetchCategory(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Name != tt.want.category.Name {
				t.Errorf("FetchCategory() got: %q, want: %q", got.Name, tt.want.category.Name)
			}
		})
	}
	if err := config.ClearCategoriesTable(config.DBDrive, dbConnStr); err != nil {
		t.Errorf("test: clear categories table: %v", err)
	}
}

func TestRepository_RemoveCategory(t *testing.T) {
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
	db, err := sql.Open(config.DBDrive, dbConnStr)
	if err != nil {
		t.Errorf("test: failed to open DB: %v\n", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("test: failed to close DB: %v\n", err)
		}
	}()
	ctx := context.Background()
	r := NewRepository(ctx, db)
	for _, g := range testdata.FakeCategories {
		err = g.InsertG(ctx, boil.Infer())
		if err != nil {
			t.Errorf("test: insert category: %s", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.RemoveCategory(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != tt.want {
				t.Errorf("RemoveCategory() got: %s, want: %q", err, tt.want)
			}
		})
	}
	if err := config.ClearCategoriesTable(config.DBDrive, dbConnStr); err != nil {
		t.Errorf("test: clear categories table: %v", err)
	}
}

func TestRepository_UpdateCategory(t *testing.T) {
	const (
		fakeExistName            = "action"
		fakeDoesNotExistName     = "fakeDoesNotExistName"
		fakeNewExistName         = "violent"
		fakeNewDoestNotExistName = "new_action"
	)
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
				fakeDoesNotExistName,
				crud.CategoryDTO{
					Name: fakeNewDoestNotExistName,
				},
			},
			want:    sql.ErrNoRows,
			wantErr: true,
		},
		{
			name: "When name in CategoryDTO already exists",
			args: args{
				fakeExistName,
				crud.CategoryDTO{
					Name: fakeNewExistName,
				},
			},
			want:    fmt.Errorf("%s %w", fakeNewExistName, logger.ErrAlreadyExists),
			wantErr: true,
		},
		{
			name: "When name exists and CategoryDTO is right",
			args: args{
				fakeExistName,
				crud.CategoryDTO{
					Name: fakeNewDoestNotExistName,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	db, err := sql.Open(config.DBDrive, dbConnStr)
	if err != nil {
		t.Errorf("test: failed to open DB: %v\n", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("test: failed to close DB: %v\n", err)
		}
	}()
	ctx := context.Background()
	r := NewRepository(ctx, db)
	for _, g := range testdata.FakeCategories {
		err = g.InsertG(ctx, boil.Infer())
		if err != nil {
			t.Errorf("test: insert category: %s", err)
			return
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := r.UpdateCategory(tt.args.name, tt.args.categoryDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("UpdateCategory() got: %v, want: %v", err, tt.want)
			}
		})
	}
	if err := config.ClearCategoriesTable(config.DBDrive, dbConnStr); err != nil {
		t.Errorf("test: clear categories table: %v", err)
	}
}

func TestCategory_isValidUUIDHook(t *testing.T) {
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
	db, err := sql.Open(config.DBDrive, dbConnStr)
	if err != nil {
		t.Errorf("test: failed to open DB: %v\n", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("test: failed to close DB: %v\n", err)
		}
	}()
	ctx := context.Background()
	r := NewRepository(ctx, db)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.category.InsertG(r.ctx, boil.Infer())
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
