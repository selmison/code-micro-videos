package crud_test

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud"
	"github.com/selmison/code-micro-videos/pkg/crud/mock"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/pkg/storage/sqlboiler"
	"github.com/selmison/code-micro-videos/testdata"
)

func TestAddGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeName := strings.ToLower(faker.FirstName())
	fakeDoesNotExistCategory := crud.CategoryDTO{Name: faker.FirstName()}
	type fields struct {
		r sqlboiler.Repository
	}
	type returns struct {
		genre models.Genre
		err   error
	}
	type args struct {
		dto crud.GenreDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    returns
		wantErr bool
	}{
		{
			name:    "When GenreDTO is not provided",
			args:    args{crud.GenreDTO{}},
			want:    returns{err: fmt.Errorf("'Name' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When the Name in GenreDTO is blank",
			args: args{crud.GenreDTO{
				Name: "    ",
			}},
			want:    returns{err: fmt.Errorf("'Name' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When the Name in GenreDTO already exists",
			args: args{crud.GenreDTO{
				Name:       strings.ToLower(faker.FirstName()),
				Categories: []crud.CategoryDTO{fakeDoesNotExistCategory},
			}},
			want:    returns{err: logger.ErrAlreadyExists},
			wantErr: true,
		},
		{
			name: "When GenreDTO is with wrong genres",
			args: args{
				crud.GenreDTO{
					Name:       fakeName,
					Categories: []crud.CategoryDTO{fakeDoesNotExistCategory},
				},
			},
			want:    returns{err: logger.ErrIsRequired},
			wantErr: true,
		},
		{
			name: "When GenreDTO is right",
			args: args{crud.GenreDTO{
				Name:       fakeName,
				Categories: []crud.CategoryDTO{fakeDoesNotExistCategory},
			}},
			want: returns{models.Genre{
				Name: fakeName,
			}, nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "When the Name in GenreDTO already exists" ||
				tt.name == "When GenreDTO is with wrong genres" ||
				tt.name == "When GenreDTO is right" {
				dto := crud.GenreDTO{
					Name:       strings.ToLower(strings.TrimSpace(tt.args.dto.Name)),
					Categories: tt.args.dto.Categories,
				}
				mockR.EXPECT().
					AddGenre(dto).
					Return(tt.want.err)
			}
			s := crud.NewService(mockR)
			err := s.AddGenre(tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddGenre() error = '%v', wantErr '%v'", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("AddGenre() got = '%v', want '%v'", err, tt.want.err)
			}
		})
	}
}

func Test_service_RemoveGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	indexRandom := rand.Intn(len(testdata.FakeCategories))
	fakeNames := [2]string{
		faker.FirstName(),
		testdata.FakeCategories[indexRandom].Name,
	}
	type fields struct {
		r sqlboiler.Repository
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
			name:    "When name is blank",
			args:    args{"     "},
			want:    fmt.Errorf("'name' %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name:    "When name is not found",
			args:    args{fakeNames[0]},
			want:    fmt.Errorf("%s: %w", fakeNames[0], logger.ErrNotFound),
			wantErr: true,
		},
		{
			name:    "When name is found",
			args:    args{fakeNames[1]},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "When name is not found" ||
				tt.name == "When name is found" {
				name := strings.ToLower(strings.TrimSpace(tt.args.name))
				mockR.EXPECT().
					RemoveGenre(name).
					Return(tt.want)
			}
			s := crud.NewService(mockR)
			err := s.RemoveGenre(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveGenre() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("RemoveGenre() got: '%v', want: '%v'", err, tt.want)
			}
		})
	}
}

func Test_service_UpdateGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	const (
		fakeExistName        = "fakeExistName"
		fakeDoesNotExistName = "fakeDoesNotExistName"
		fakeCategoryIndex    = 0
	)
	fakeName := faker.FirstName()
	fakeDoesNotExistCategoryDTO := crud.CategoryDTO{Name: faker.FirstName()}
	fakeExistCategoryDTO := testdata.FakeCategoriesDTO[fakeCategoryIndex]
	type fields struct {
		r sqlboiler.Repository
	}
	type args struct {
		name string
		dto  crud.GenreDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "When name is blank",
			args: args{
				"     ",
				crud.GenreDTO{
					Name: faker.FirstName(),
				},
			},
			want:    fmt.Errorf("'name' %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name: "When the Name in CategoryDTO is blank",
			args: args{
				fakeExistName,
				crud.GenreDTO{
					Name: "    ",
				}},
			want:    fmt.Errorf("'Name' field %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name: "When GenreDTO is with wrong genres",
			args: args{
				fakeExistName,
				crud.GenreDTO{
					Name:       fakeName,
					Categories: []crud.CategoryDTO{fakeDoesNotExistCategoryDTO},
				},
			},
			want:    logger.ErrIsRequired,
			wantErr: true,
		},
		{
			name:    "When GenreDTO is not provided",
			args:    args{fakeExistName, crud.GenreDTO{}},
			want:    fmt.Errorf("'Name' field %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name: "When name is not found",
			args: args{
				fakeDoesNotExistName,
				crud.GenreDTO{
					Name:       faker.FirstName(),
					Categories: []crud.CategoryDTO{fakeExistCategoryDTO},
				},
			},
			want:    fmt.Errorf("%s: %w", fakeDoesNotExistName, logger.ErrNotFound),
			wantErr: true,
		},
		{
			name: "When GenreDTO is right",
			args: args{
				fakeExistName,
				crud.GenreDTO{
					Name:       faker.FirstName(),
					Categories: []crud.CategoryDTO{fakeDoesNotExistCategoryDTO},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "When the Name in GenreDTO already exists" ||
				tt.name == "When GenreDTO is with wrong genres" ||
				tt.name == "When name is not found" ||
				tt.name == "When GenreDTO is right" {
				name := strings.ToLower(strings.TrimSpace(tt.args.name))
				dto := crud.GenreDTO{
					Name:       strings.ToLower(strings.TrimSpace(tt.args.dto.Name)),
					Categories: tt.args.dto.Categories,
				}
				mockR.EXPECT().
					UpdateGenre(name, dto).
					Return(tt.want)
			}
			s := crud.NewService(mockR)
			err := s.UpdateGenre(tt.args.name, tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateGenre() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.want.Error() {
				t.Errorf("AddCategory() got = '%v', want '%v'", err, tt.want)
			}
		})
	}
}

func Test_service_GetGenres(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeGenreSlice := models.GenreSlice{
		&models.Genre{
			Name: "action",
		},
		&models.Genre{
			Name: "fiction",
		},
		&models.Genre{
			Name: "animation",
		},
	}
	fakeLimit := len(fakeGenreSlice)
	type args struct {
		limit int
	}
	type returns struct {
		cs models.GenreSlice
		e  error
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
			want:    returns{fakeGenreSlice, nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockR.EXPECT().
					GetGenres(tt.args.limit).
					Return(
						fakeGenreSlice,
						nil,
					)
			}
			s := crud.NewService(mockR)
			got, err := s.GetGenres(tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGenres() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.cs) {
				t.Errorf("GetGenres() got = %v, want %v", got, tt.want.cs)
			}
			if !reflect.DeepEqual(err, tt.want.e) {
				t.Errorf("GetGenres() got = %v, want %v", err, tt.want.e)
			}
		})
	}
}

func Test_service_FetchGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeDoesNotExistName := "fakeDoesNotExistName"
	fakeExistName := "action"
	fakeErrorInternalApplication := fmt.Errorf("Service.FetchGenre(): %w", logger.ErrInternalApplication)
	type args struct {
		name string
	}
	type returns struct {
		c models.Genre
		e error
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
				models.Genre{},
				fakeErrorInternalApplication,
			},
			wantErr: true,
			setupMockR: func() {
				mockR.EXPECT().
					FetchGenre("anyname").
					Return(
						models.Genre{},
						fakeErrorInternalApplication,
					)
			},
		},
		{
			name: "When name is not found",
			args: args{fakeDoesNotExistName},
			want: returns{
				models.Genre{},
				fmt.Errorf("%s: %w", fakeDoesNotExistName, logger.ErrNotFound),
			},
			wantErr: true,
			setupMockR: func() {
				mockR.EXPECT().
					FetchGenre(strings.ToLower(strings.TrimSpace(fakeDoesNotExistName))).
					Return(
						models.Genre{},
						sql.ErrNoRows,
					)
			},
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
			setupMockR: func() {
				mockR.EXPECT().
					FetchGenre(fakeExistName).
					Return(
						models.Genre{
							Name: fakeExistName,
						},
						nil,
					)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMockR()
			s := crud.NewService(mockR)
			got, err := s.FetchGenre(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchGenre() error: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.c) {
				t.Errorf("GetGenre() got: %v, want: %v", got, tt.want.c)
			}
			if tt.wantErr && errors.Is(err, tt.want.e) {
				t.Errorf("GetGenre() got: %v, want: %v", err, tt.want.e)
			}
		})
	}
}
