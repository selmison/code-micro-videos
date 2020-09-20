package service_test

import (
	"context"
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
	"github.com/selmison/code-micro-videos/pkg/crud/domain"
	"github.com/selmison/code-micro-videos/pkg/crud/service/mock"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/pkg/storage/sqlboiler"
	"github.com/selmison/code-micro-videos/testdata"
)

func TestCreateCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	const fakeGenreIndex = 0
	fakeCtx := context.Background()
	fakeName := faker.FirstName()
	fakeDescription := faker.Sentence()
	fakeDoesNotExistGenre := domain.Genre{Name: faker.FirstName()}
	fakeExistGenreDTO := testdata.FakeGenresDTO[fakeGenreIndex]
	fakeExistGenreDTOs := []domain.Genre{
		{fakeExistGenreDTO.Id, fakeDoesNotExistGenre.Name},
	}
	type fields struct {
		r sqlboiler.Repository
	}
	type returns struct {
		category models.Category
		err      error
	}
	type args struct {
		ctx context.Context
		dto domain.Category
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    returns
		wantErr bool
	}{
		{
			name:    "When Category is not provided",
			args:    args{fakeCtx, domain.Category{}},
			want:    returns{err: fmt.Errorf("genres %w", logger.ErrIsEmpty)},
			wantErr: true,
		},
		{
			name: "When the Name in Category is blank",
			args: args{fakeCtx, domain.Category{
				Name:        "    ",
				Description: fakeDescription,
			}},
			want:    returns{err: fmt.Errorf("'name' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When the Name in Category already exists",
			args: args{fakeCtx, domain.Category{
				Name:        fakeName,
				Description: fakeDescription,
				Genres:      fakeExistGenreDTOs,
			},
			},
			want:    returns{err: logger.ErrAlreadyExists},
			wantErr: true,
		},
		{
			name: "When Category is with wrong genres",
			args: args{fakeCtx,
				domain.Category{
					Name:        fakeName,
					Description: fakeDescription,
					Genres: []domain.Genre{
						{fakeDoesNotExistGenre.Id, fakeDoesNotExistGenre.Name},
					},
				},
			},
			want:    returns{err: logger.ErrIsRequired},
			wantErr: true,
		},
		{
			name: "When Category is right",
			args: args{fakeCtx, domain.Category{
				Name:        fakeName,
				Description: fakeDescription,
				Genres:      fakeExistGenreDTOs,
			},
			},
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
			if tt.name == "When the Name in Category already exists" ||
				tt.name == "When Category is with wrong genres" ||
				tt.name == "When Category is right" {
				desc := strings.TrimSpace(tt.args.dto.Description)
				category := domain.Category{
					Name:        strings.ToLower(strings.TrimSpace(tt.args.dto.Name)),
					Description: desc,
					Genres:      tt.args.dto.Genres,
				}
				mockR.EXPECT().
					CreateCategory(tt.args.ctx, category).
					Return(tt.want.err)
			}
			s := domain.NewService(mockR)
			err := s.CreateCategory(tt.args.ctx, tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateCategory() error = '%v', wantErr '%v'", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("CreateCategory() got = '%v', want '%v'", err, tt.want.err)
			}
		})
	}
}

func Test_service_RemoveCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeCtx := context.Background()
	indexRandom := rand.Intn(len(testdata.FakeCategories))
	fakeNames := [2]string{
		faker.FirstName(),
		testdata.FakeCategories[indexRandom].Name,
	}
	type fields struct {
		r sqlboiler.Repository
	}
	type args struct {
		ctx  context.Context
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
			args:    args{fakeCtx, "     "},
			want:    fmt.Errorf("'name' %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name:    "When name is not found",
			args:    args{fakeCtx, fakeNames[0]},
			want:    fmt.Errorf("%s: %w", fakeNames[0], logger.ErrNotFound),
			wantErr: true,
		},
		{
			name:    "When name is found",
			args:    args{fakeCtx, fakeNames[1]},
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
					RemoveCategory(tt.args.ctx, name).
					Return(tt.want)
			}
			s := domain.NewService(mockR)
			err := s.RemoveCategory(tt.args.ctx, tt.args.name)
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
		fakeGenreIndex       = 0
	)
	fakeCtx := context.Background()
	fakeName := faker.FirstName()
	fakeDescription := faker.Sentence()
	fakeDoesNotExistGenre := domain.Genre{Name: faker.FirstName()}
	fakeExistGenreDTO := testdata.FakeGenresDTO[fakeGenreIndex]
	fakeExistGenreDTOs := []domain.Genre{
		{fakeExistGenreDTO.Id, fakeExistGenreDTO.Name},
	}
	type fields struct {
		r sqlboiler.Repository
	}
	type args struct {
		ctx  context.Context
		name string
		dto  domain.Category
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
			args: args{fakeCtx,
				"     ",
				domain.Category{
					Name:        faker.FirstName(),
					Description: &fakeDescription,
				},
			},
			want:    fmt.Errorf("'name' %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name: "When the Name in Category is blank",
			args: args{fakeCtx,
				fakeExistName,
				domain.Category{
					Name:        "    ",
					Description: fakeDescription,
				}},
			want:    fmt.Errorf("'name' field %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name: "When Category is with wrong genres",
			args: args{fakeCtx,
				fakeExistName,
				domain.Category{
					Name:        fakeName,
					Description: fakeDescription,
					Genres: []domain.Genre{
						{fakeDoesNotExistGenre.Id, fakeDoesNotExistGenre.Name},
					},
				},
			},
			want:    logger.ErrIsRequired,
			wantErr: true,
		},
		{
			name:    "When Category is not provided",
			args:    args{fakeCtx, fakeExistName, domain.Category{}},
			want:    fmt.Errorf("category %w", logger.ErrIsEmpty),
			wantErr: true,
		},
		{
			name: "When name is not found",
			args: args{fakeCtx,
				fakeDoesNotExistName,
				domain.Category{
					Name:        faker.FirstName(),
					Description: fakeDescription,
					Genres:      fakeExistGenreDTOs,
				},
			},
			want:    fmt.Errorf("%s: %w", fakeDoesNotExistName, logger.ErrNotFound),
			wantErr: true,
		},
		{
			name: "When Category is right",
			args: args{fakeCtx,
				fakeExistName,
				domain.Category{
					Name:        faker.FirstName(),
					Description: fakeDescription,
					Genres:      fakeExistGenreDTOs,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "When the Name in Category already exists" ||
				tt.name == "When Category is with wrong genres" ||
				tt.name == "When name is not found" ||
				tt.name == "When Category is right" {
				name := strings.ToLower(strings.TrimSpace(tt.args.name))
				category := domain.Category{
					Name:        strings.ToLower(strings.TrimSpace(tt.args.dto.Name)),
					Description: tt.args.dto.Description,
					Genres:      tt.args.dto.Genres,
				}
				mockR.EXPECT().
					UpdateCategory(tt.args.ctx, name, category).
					Return(tt.want)
			}
			s := domain.NewService(mockR)
			err := s.UpdateCategory(tt.args.ctx, tt.args.name, tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCategory() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.want.Error() {
				t.Errorf("UpdateCategory() got = '%v', want '%v'", err, tt.want)
			}
		})
	}
}

func Test_service_GetCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeCtx := context.Background()
	fakeDesc := faker.Sentence()
	fakeCategory := []domain.Category{
		{
			Name: faker.FirstName(),
		},
		{
			Name:        faker.FirstName(),
			Description: &fakeDesc,
		},
		{
			Name: faker.FirstName(),
		},
	}
	fakeLimit := len(fakeCategory)
	type args struct {
		ctx   context.Context
		limit int
	}
	type returns struct {
		categories []domain.Category
		err        error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name:    "When limit is less than zero",
			args:    args{fakeCtx, -1},
			want:    returns{nil, logger.ErrInvalidedLimit},
			wantErr: true,
		},
		{
			name:    "When limit is right",
			args:    args{fakeCtx, fakeLimit},
			want:    returns{fakeCategory, nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockR.EXPECT().
					GetCategories(tt.args.ctx, tt.args.limit).
					Return(
						fakeCategory,
						nil,
					)
			}
			s := domain.NewService(mockR)
			got, err := s.GetCategories(tt.args.ctx, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.categories) {
				t.Errorf("GetCategories() got = %v, want %v", got, tt.want.categories)
			}
			if !reflect.DeepEqual(err, tt.want.err) {
				t.Errorf("GetCategories() got = %v, want %v", err, tt.want.err)
			}
		})
	}
}

func Test_service_FetchCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeCtx := context.Background()
	indexRandom := rand.Intn(len(testdata.FakeCategories))
	fakeExistsName := testdata.FakeCategories[indexRandom].Name
	fakeDoesNotExistName := faker.FirstName()
	type args struct {
		ctx  context.Context
		name string
	}
	type returns struct {
		category domain.Category
		err      error
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
				domain.Category{},
				fmt.Errorf("%s: %w", fakeDoesNotExistName, logger.ErrNotFound),
			},
			wantErr: true,
		},
		{
			name: "When name is found",
			args: args{fakeCtx, fakeExistsName},
			want: returns{
				domain.Category{
					Name: fakeExistsName,
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
					FetchCategory(tt.args.ctx, tt.args.name).
					Return(
						tt.want.category,
						nil,
					)
			} else {
				mockR.EXPECT().
					FetchCategory(tt.args.ctx, strings.ToLower(tt.args.name)).
					Return(
						domain.Category{},
						sql.ErrNoRows,
					)
			}
			s := domain.NewService(mockR)
			got, err := s.FetchCategory(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.category) {
				t.Errorf("GetCategory() got = %v, want: %v", got, tt.want.category)
			}
			if tt.wantErr && errors.Is(err, tt.want.err) {
				t.Errorf("GetCategory() got: %v, want: %v", err, tt.want.err)
			}
		})
	}
}
