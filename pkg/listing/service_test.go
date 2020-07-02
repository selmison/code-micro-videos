package listing_test

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/listing"
	"github.com/selmison/code-micro-videos/pkg/listing/mock"
	"github.com/selmison/code-micro-videos/testdata/seeds"
	"github.com/volatiletech/null/v8"
	"math/rand"
	"reflect"
	"testing"
)

var (
	seedsArray []seeds.Seed
)

func init() {
	faker.SetGenerateUniqueValues(true)
	seedsArray = seeds.MakeSeeds(10)
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
			want:    returns{nil, listing.ErrInvalidedLimit},
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
			s := listing.NewService(mockR)
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
				fmt.Errorf("%s: %w", fakeNames[0], listing.ErrNotFound),
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
					FetchCategory(tt.args.name).
					Return(
						models.Category{},
						sql.ErrNoRows,
					)
			}
			s := listing.NewService(mockR)
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
