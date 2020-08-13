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
	"github.com/volatiletech/null/v8"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud"
	"github.com/selmison/code-micro-videos/pkg/crud/mock"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/pkg/storage/sqlboiler"
	"github.com/selmison/code-micro-videos/testdata/seeds"
)

var (
	seedsArray []seeds.Seed
)

func init() {
	faker.SetGenerateUniqueValues(true)
	seedsArray = seeds.MakeSeeds(10)
}

func TestAddCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeName := strings.ToLower(faker.FirstName())
	fakeDescription := faker.Sentence()
	type fields struct {
		r sqlboiler.Repository
	}
	type returns struct {
		c models.Category
		e error
	}
	type args struct {
		dto crud.CategoryDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    returns
		wantErr bool
	}{
		{
			name:    "When CategoryDTO is not provided",
			args:    args{crud.CategoryDTO{}},
			want:    returns{models.Category{}, logger.ErrIsRequired},
			wantErr: true,
		},
		{
			name: "When the name in CategoryDTO is blank",
			args: args{crud.CategoryDTO{
				Name:        "    ",
				Description: faker.Sentence(),
			}},
			want:    returns{models.Category{}, logger.ErrIsRequired},
			wantErr: true,
		},
		{
			name: "When name in CategoryDTO already exists",
			args: args{crud.CategoryDTO{
				Name: strings.ToLower(faker.FirstName()),
			}},
			want:    returns{models.Category{}, logger.ErrAlreadyExists},
			wantErr: true,
		},
		{
			name: "When CategoryDTO is right",
			args: args{crud.CategoryDTO{
				Name:        fakeName,
				Description: fakeDescription,
			}},
			want: returns{models.Category{
				Name: fakeName,
				Description: null.String{
					String: fakeDescription},
			}, nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want.e != logger.ErrIsRequired {
				mockR.EXPECT().
					AddCategory(tt.args.dto).
					Return(tt.want.e)
			}
			s := crud.NewService(mockR)
			err := s.AddCategory(tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCategory() error = '%v', wantErr '%v'", err, tt.wantErr)
				return
			}
			if !errors.Is(err, tt.want.e) {
				t.Errorf("AddCategory() got = '%v', want '%v'", err, tt.want.e)
			}
		})
	}
}

func Test_service_RemoveCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	indexRandom := rand.Intn(len(seedsArray))
	fakeNames := [2]string{
		faker.FirstName(),
		seedsArray[indexRandom].Name,
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
			if tt.name == "When name is not found" {
				mockR.EXPECT().
					RemoveCategory(tt.args.name).
					Return(tt.want)
			} else if tt.name == "When name is found" {
				mockR.EXPECT().
					RemoveCategory(tt.args.name).
					Return(tt.want)
			}
			s := crud.NewService(mockR)
			err := s.RemoveCategory(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveCategory() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("RemoveCategory() got: '%v', want: '%v'", err, tt.want)
			}
		})
	}
}

func Test_service_UpdateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	const (
		fakeExistName        = "fakeExistName"
		fakeDoesNotExistName = "fakeDoesNotExistName"
	)
	type fields struct {
		r sqlboiler.Repository
	}
	type args struct {
		name string
		dto  crud.CategoryDTO
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
				crud.CategoryDTO{
					Name:        faker.FirstName(),
					Description: faker.Sentence(),
				},
			},
			want:    logger.ErrIsRequired,
			wantErr: true,
		},
		{
			name:    "When CategoryDTO is not provided",
			args:    args{fakeExistName, crud.CategoryDTO{}},
			want:    logger.ErrIsRequired,
			wantErr: true,
		},
		{
			name: "When name is not found",
			args: args{
				fakeDoesNotExistName,
				crud.CategoryDTO{
					Name:        faker.FirstName(),
					Description: faker.Sentence(),
				},
			},
			want:    fmt.Errorf("%s: %w", fakeDoesNotExistName, logger.ErrNotFound),
			wantErr: true,
		},
		{
			name: "When name is found and CategoryDTO is provided",
			args: args{
				fakeExistName,
				crud.CategoryDTO{
					Name:        faker.FirstName(),
					Description: faker.Sentence(),
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "When name is not found" || tt.name == "When name is found and CategoryDTO is provided" {
				mockR.EXPECT().
					UpdateCategory(tt.args.name, tt.args.dto).
					Return(tt.want)
			}
			s := crud.NewService(mockR)
			err := s.UpdateCategory(tt.args.name, tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCategory() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !errors.Is(err, tt.want) {
				t.Errorf("UpdateCategory() got: '%v', want: '%v'", err, tt.want)
			}
		})
	}
}

func Test_service_GetCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeCategorySlice := models.CategorySlice{
		&models.Category{
			Name:        faker.FirstName(),
			Description: null.String{String: seeds.Sentence(), Valid: true},
		},
		&models.Category{
			Name:        faker.FirstName(),
			Description: null.String{String: seeds.Sentence(), Valid: true},
		},
		&models.Category{
			Name:        faker.FirstName(),
			Description: null.String{String: seeds.Sentence(), Valid: true},
		},
	}
	fakeLimit := len(fakeCategorySlice)
	type args struct {
		limit int
	}
	type returns struct {
		cs models.CategorySlice
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
			want:    returns{fakeCategorySlice, nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockR.EXPECT().
					GetCategories(tt.args.limit).
					Return(
						fakeCategorySlice,
						nil,
					)
			}
			s := crud.NewService(mockR)
			got, err := s.GetCategories(tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.cs) {
				t.Errorf("GetCategories() got = %v, want %v", got, tt.want.cs)
			}
			if !reflect.DeepEqual(err, tt.want.e) {
				t.Errorf("GetCategories() got = %v, want %v", err, tt.want.e)
			}
		})
	}
}

func Test_service_FetchCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	indexRandom := rand.Intn(len(seedsArray))
	fakeNames := [2]string{
		faker.FirstName(),
		seedsArray[indexRandom].Name,
	}
	type args struct {
		name string
	}
	type returns struct {
		c models.Category
		e error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "When name is not found",
			args: args{fakeNames[0]},
			want: returns{
				models.Category{},
				fmt.Errorf("%s: %w", fakeNames[0], logger.ErrNotFound),
			},
			wantErr: true,
		},
		{
			name: "When name is found",
			args: args{fakeNames[1]},
			want: returns{
				models.Category{
					Name: fakeNames[1],
				},
				nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockR.EXPECT().
					FetchCategory(tt.args.name).
					Return(
						tt.want.c,
						nil,
					)
			} else {
				mockR.EXPECT().
					FetchCategory(strings.ToLower(tt.args.name)).
					Return(
						models.Category{},
						sql.ErrNoRows,
					)
			}
			s := crud.NewService(mockR)
			got, err := s.FetchCategory(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.c) {
				t.Errorf("GetCategory() got = %v, want: %v", got, tt.want.c)
			}
			if tt.wantErr && errors.Is(err, tt.want.e) {
				t.Errorf("GetCategory() got: %v, want: %v", err, tt.want.e)
			}
		})
	}
}
