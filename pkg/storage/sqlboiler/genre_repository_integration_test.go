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

func TestRepository_AddGenre(t *testing.T) {
	cfg, teardownTestCase, repository, err := setupTestCase(nil)
	if err != nil {
		t.Errorf("test: failed to open DB: %v\n", err)
	}
	defer teardownTestCase(t)
	const (
		fakeExistName        = "action"
		fakeDoesNotExistName = "fakeDoesNotExistName"
	)
	fakeExistGenre := models.Genre{
		ID:   uuid.New().String(),
		Name: fakeExistName,
	}
	type args struct {
		genreDTO crud.GenreDTO
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
			name: "When name in GenreDTO already exists",
			args: args{crud.GenreDTO{
				Name: fakeExistName,
			}},
			want:    returns{models.Genre{}, fmt.Errorf("name '%s' %w", fakeExistName, logger.ErrAlreadyExists)},
			wantErr: true,
		},
		{
			name: "When GenreDTO is right",
			args: args{
				crud.GenreDTO{
					Name: fakeDoesNotExistName,
				},
			},
			want: returns{
				models.Genre{
					Name: fakeDoesNotExistName,
				},
				nil,
			},
			wantErr: false,
		},
	}
	db, err := sql.Open(cfg.DBDrive, cfg.DBConnStr)
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
	err = fakeExistGenre.InsertG(ctx, boil.Infer())
	if err != nil {
		t.Errorf("test: insert genre: %s", err)
		return
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.AddGenre(tt.args.genreDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddGenre() error: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(err, tt.want.err) {
				t.Errorf("AddGenre() got: %v, want: %v", err, tt.want.err)
				return
			}
		})
	}
}

func TestRepository_GetGenres(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeGenres)
	if err != nil {
		t.Errorf("test: failed to open DB: %v\n", err)
	}
	defer teardownTestCase(t)
	maximum := len(testdata.FakeGenres)
	type args struct {
		limit int
	}
	type returns struct {
		genres models.GenreSlice
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
			got, err := repository.GetGenres(tt.args.limit)
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
		t.Errorf("test: failed to open DB: %v\n", err)
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
			args: args{fakeDoesNotExistName},
			want: returns{
				models.Genre{},
				sql.ErrNoRows,
			},
			wantErr: true,
		},
		{
			name: "When name is found",
			args: args{fakeExistName},
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
			got, err := repository.FetchGenre(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchGenre() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Name != tt.want.genre.Name {
				t.Errorf("FetchGenre() got: %q, want: %q", got.Name, tt.want.genre.Name)
			}
		})
	}
}

func TestRepository_RemoveGenre(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeGenres)
	if err != nil {
		t.Errorf("test: failed to open DB: %v\n", err)
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
			err := repository.RemoveGenre(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveGenre() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != tt.want {
				t.Errorf("RemoveGenre() got: %s, want: %q", err, tt.want)
			}
		})
	}
}

func TestRepository_UpdateGenre(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeGenres)
	if err != nil {
		t.Errorf("test: failed to open DB: %v\n", err)
	}
	defer teardownTestCase(t)
	const (
		fakeExistName            = "action"
		fakeDoesNotExistName     = "fakeDoesNotExistName"
		fakeNewExistName         = "violent"
		fakeNewDoestNotExistName = "new_action"
	)
	type fields struct {
		ctx context.Context
	}
	type args struct {
		name     string
		genreDTO crud.GenreDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "When name to update doesn't exist",
			args: args{
				fakeDoesNotExistName,
				crud.GenreDTO{
					Name: fakeNewDoestNotExistName,
				},
			},
			want:    sql.ErrNoRows,
			wantErr: true,
		},
		{
			name: "When name in GenreDTO already exists",
			args: args{
				fakeExistName,
				crud.GenreDTO{
					Name: fakeNewExistName,
				},
			},
			want:    fmt.Errorf("%s %w", fakeNewExistName, logger.ErrAlreadyExists),
			wantErr: true,
		},
		{
			name: "When name exists and GenreDTO is right",
			args: args{
				fakeExistName,
				crud.GenreDTO{
					Name: fakeNewDoestNotExistName,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.UpdateGenre(tt.args.name, tt.args.genreDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateGenre() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("UpdateGenre() got: %v, want: %v", err, tt.want)
			}
		})
	}
}

//func ClearGenresTable(cfg.DBConnStr string) error {
//	db, err := sql.Open(cfg.DBDrive, cfg.DBConnStr)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		if err := db.Close(); err != nil {
//			log.Fatalln(err)
//		}
//	}()
//	if _, err := models.Genres().DeleteAll(context.Background(), db, true); err != nil {
//		return err
//	}
//	return nil
//}

func TestGenre_isValidUUIDHook(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(nil)
	if err != nil {
		t.Errorf("test: failed to open DB: %v\n", err)
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
			name: "When UUID is not validated",
			args: args{
				models.Genre{
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
				t.Errorf("isValidUUIDHook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("UpdateCategory() got: %v, want: %v", err, tt.want)
			}
		})
	}
}
