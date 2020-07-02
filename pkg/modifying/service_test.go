package modifying_test

import (
	"errors"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/listing"
	"github.com/selmison/code-micro-videos/pkg/modifying"
	"github.com/selmison/code-micro-videos/pkg/modifying/mock"
	"github.com/selmison/code-micro-videos/pkg/storage/sqlboiler"
	"github.com/selmison/code-micro-videos/testdata/seeds"
	"github.com/volatiletech/null/v8"
	"math/rand"
	"reflect"
	"strings"
	"testing"
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
		c modifying.CategoryDTO
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
			args:    args{modifying.CategoryDTO{}},
			want:    returns{models.Category{}, modifying.ErrIsRequired},
			wantErr: true,
		},
		{
			name: "When the name in CategoryDTO is blank",
			args: args{modifying.CategoryDTO{
				Name:        "    ",
				Description: faker.Sentence(),
			}},
			want:    returns{models.Category{}, modifying.ErrIsRequired},
			wantErr: true,
		},
		{
			name: "When name in CategoryDTO already exists",
			args: args{modifying.CategoryDTO{
				Name: strings.ToLower(faker.FirstName()),
			}},
			want:    returns{models.Category{}, modifying.ErrAlreadyExists},
			wantErr: true,
		},
		{
			name: "When CategoryDTO is right",
			args: args{modifying.CategoryDTO{
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
			if tt.want.e != modifying.ErrIsRequired {
				mockR.EXPECT().
					AddCategory(tt.args.c).
					Return(tt.want.e)
			}
			s := modifying.NewService(mockR)
			err := s.AddCategory(tt.args.c)
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
			want:    fmt.Errorf("'name' %w", modifying.ErrIsRequired),
			wantErr: true,
		},
		{
			name:    "When name is not found",
			args:    args{fakeNames[0]},
			want:    fmt.Errorf("%s: %w", fakeNames[0], listing.ErrNotFound),
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
			s := modifying.NewService(mockR)
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
	indexRandom := rand.Intn(len(seedsArray))
	fakeNames := [2]string{
		faker.FirstName(),
		seedsArray[indexRandom].Name,
	}
	//fakeDescription := faker.Sentence()
	type fields struct {
		r sqlboiler.Repository
	}
	type args struct {
		name string
		c    modifying.CategoryDTO
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
				modifying.CategoryDTO{
					Name:        faker.FirstName(),
					Description: faker.Sentence(),
				},
			},
			want:    modifying.ErrIsRequired,
			wantErr: true,
		},
		{
			name:    "When CategoryDTO is not provided",
			args:    args{fakeNames[0], modifying.CategoryDTO{}},
			want:    modifying.ErrIsRequired,
			wantErr: true,
		},
		{
			name: "When name is not found",
			args: args{
				fakeNames[0],
				modifying.CategoryDTO{
					Name:        faker.FirstName(),
					Description: faker.Sentence(),
				},
			},
			want:    fmt.Errorf("%s: %w", fakeNames[0], listing.ErrNotFound),
			wantErr: true,
		},
		{
			name: "When name is found and CategoryDTO is provided",
			args: args{
				fakeNames[1],
				modifying.CategoryDTO{
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
			if tt.name == "When name is not found" {
				mockR.EXPECT().
					UpdateCategory(tt.args.name, tt.args.c).
					Return(tt.want)
			} else if tt.name == "When name is found and CategoryDTO is provided" {
				mockR.EXPECT().
					UpdateCategory(tt.args.name, tt.args.c).
					Return(tt.want)
			}
			s := modifying.NewService(mockR)
			err := s.UpdateCategory(tt.args.name, tt.args.c)
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
