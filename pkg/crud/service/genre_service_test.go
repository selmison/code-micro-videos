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

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud/domain"
	"github.com/selmison/code-micro-videos/pkg/crud/service"
	"github.com/selmison/code-micro-videos/pkg/crud/service/mock"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/pkg/storage/sqlboiler"
	"github.com/selmison/code-micro-videos/testdata"
)

func TestCreateGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeCtx := context.Background()
	fakeName := strings.ToLower(faker.FirstName())
	fakeDoesNotExistCategory := domain.CategoryValidatable{Name: faker.FirstName()}
	fakeDoesNotExistCategories := []domain.CategoryValidatable{
		{
			Id:          fakeDoesNotExistCategory.Id,
			Name:        fakeDoesNotExistCategory.Name,
			Description: fakeDoesNotExistCategory.Description,
		},
	}
	type fields struct {
		r sqlboiler.Repository
	}
	type returns struct {
		genre models.Genre
		err   error
	}
	type args struct {
		ctx   context.Context
		genre domain.GenreValidatable
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    returns
		wantErr bool
	}{
		{
			name:    "When Genre is not provided",
			args:    args{fakeCtx, domain.GenreValidatable{}},
			want:    returns{err: fmt.Errorf("genres %w", logger.ErrIsEmpty)},
			wantErr: true,
		},
		{
			name: "When the Name in Genre is blank",
			args: args{fakeCtx, domain.GenreValidatable{
				Name: "    ",
			}},
			want:    returns{err: fmt.Errorf("'name' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When the Name in Genre already exists",
			args: args{fakeCtx, domain.GenreValidatable{
				Name:       strings.ToLower(faker.FirstName()),
				Categories: fakeDoesNotExistCategories,
			}},
			want:    returns{err: logger.ErrAlreadyExists},
			wantErr: true,
		},
		{
			name: "When Genre is with wrong genres",
			args: args{fakeCtx,
				domain.GenreValidatable{
					Name:       fakeName,
					Categories: fakeDoesNotExistCategories,
				},
			},
			want:    returns{err: logger.ErrIsRequired},
			wantErr: true,
		},
		{
			name: "When Genre is right",
			args: args{fakeCtx, domain.GenreValidatable{
				Name:       fakeName,
				Categories: fakeDoesNotExistCategories,
			},
			},
			want: returns{models.Genre{
				Name: fakeName,
			}, nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "When the Name in Genre already exists" ||
				tt.name == "When Genre is with wrong genres" ||
				tt.name == "When Genre is right" {
				genre := domain.GenreValidatable{
					Name:       strings.ToLower(strings.TrimSpace(tt.args.genre.Name)),
					Categories: tt.args.genre.Categories,
				}
				mockR.EXPECT().
					CreateGenre(tt.args.ctx, genre).
					Return(tt.want.err)
			}
			s := service.NewService(mockR)
			err := s.CreateGenre(tt.args.ctx, tt.args.genre)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateGenre() error = '%v', wantErr '%v'", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("CreateGenre() got = '%v', want '%v'", err, tt.want.err)
			}
		})
	}
}

func Test_service_RemoveGenre(t *testing.T) {
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
					RemoveGenre(tt.args.ctx, name).
					Return(tt.want)
			}
			s := service.NewService(mockR)
			err := s.RemoveGenre(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveGenreByName() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("RemoveGenreByName() got: '%v', want: '%v'", err, tt.want)
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
	fakeCtx := context.Background()
	fakeName := faker.FirstName()
	fakeDoesNotExistCategory := domain.CategoryValidatable{Name: faker.FirstName()}
	fakeExistCategory := testdata.FakeCategoriesDTO[fakeCategoryIndex]
	fakeDoesNotExistCategories := []domain.CategoryValidatable{
		{
			Id:          fakeDoesNotExistCategory.Id,
			Name:        fakeDoesNotExistCategory.Name,
			Description: fakeDoesNotExistCategory.Description,
		},
	}
	fakeExistCategories := []domain.CategoryValidatable{
		{
			Id:          fakeExistCategory.Id,
			Name:        fakeExistCategory.Name,
			Description: fakeExistCategory.Description,
		},
	}
	type fields struct {
		r sqlboiler.Repository
	}
	type args struct {
		ctx   context.Context
		name  string
		genre domain.GenreValidatable
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
				fakeCtx,
				"     ",
				domain.GenreValidatable{
					Name: faker.FirstName(),
				},
			},
			want:    fmt.Errorf("'name' %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name: "When the Name in Category is blank",
			args: args{
				fakeCtx,
				fakeExistName,
				domain.GenreValidatable{
					Name: "    ",
				}},
			want:    fmt.Errorf("'name' field %w", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name: "When Genre is with wrong genres",
			args: args{
				fakeCtx,
				fakeExistName,
				domain.GenreValidatable{
					Name:       fakeName,
					Categories: fakeDoesNotExistCategories,
				},
			},
			want:    logger.ErrIsRequired,
			wantErr: true,
		},
		{
			name: "When Genre is not provided",
			args: args{
				fakeCtx, fakeExistName, domain.GenreValidatable{}},
			want:    fmt.Errorf("genres %w", logger.ErrIsEmpty),
			wantErr: true,
		},
		{
			name: "When name is not found",
			args: args{
				fakeCtx,
				fakeDoesNotExistName,
				domain.GenreValidatable{
					Name:       faker.FirstName(),
					Categories: fakeExistCategories,
				},
			},
			want:    fmt.Errorf("%s: %w", fakeDoesNotExistName, logger.ErrNotFound),
			wantErr: true,
		},
		{
			name: "When Genre is right",
			args: args{
				fakeCtx,
				fakeExistName,
				domain.GenreValidatable{
					Name:       faker.FirstName(),
					Categories: fakeDoesNotExistCategories,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "When the Name in Genre already exists" ||
				tt.name == "When Genre is with wrong genres" ||
				tt.name == "When name is not found" ||
				tt.name == "When Genre is right" {
				name := strings.ToLower(strings.TrimSpace(tt.args.name))
				genre := domain.GenreValidatable{
					Name:       strings.ToLower(strings.TrimSpace(tt.args.genre.Name)),
					Categories: tt.args.genre.Categories,
				}
				mockR.EXPECT().
					UpdateGenre(tt.args.ctx, name, genre).
					Return(tt.want)
			}
			s := service.NewService(mockR)
			err := s.UpdateGenre(tt.args.ctx, tt.args.name, tt.args.genre)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateGenre() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.want.Error() {
				t.Errorf("CreateCategory() got = '%v', want '%v'", err, tt.want)
			}
		})
	}
}

func Test_service_GetGenres(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeCtx := context.Background()
	fakeGenres := []domain.GenreValidatable{
		{
			Name: "action",
		},
		{
			Name: "fiction",
		},
		{
			Name: "animation",
		},
	}
	fakeLimit := len(fakeGenres)
	type args struct {
		ctx   context.Context
		limit int
	}
	type returns struct {
		genres []domain.GenreValidatable
		err    error
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
			want:    returns{fakeGenres, nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockR.EXPECT().
					GetGenres(tt.args.ctx, tt.args.limit).
					Return(
						fakeGenres,
						nil,
					)
			}
			s := service.NewService(mockR)
			got, err := s.GetGenres(tt.args.ctx, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGenres() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.genres) {
				t.Errorf("GetGenres() got = %v, want %v", got, tt.want.genres)
			}
			if !reflect.DeepEqual(err, tt.want.err) {
				t.Errorf("GetGenres() got = %v, want %v", err, tt.want.err)
			}
		})
	}
}

func Test_service_FetchGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeCtx := context.Background()
	fakeDoesNotExistName := "fakeDoesNotExistName"
	fakeExistName := "action"
	fakeErrorInternalApplication := fmt.Errorf("Service.FetchGenre(): %w", logger.ErrInternalApplication)
	type args struct {
		ctx  context.Context
		name string
	}
	type returns struct {
		genre domain.GenreValidatable
		err   error
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
			args: args{fakeCtx, "anyName"},
			want: returns{
				domain.GenreValidatable{},
				fakeErrorInternalApplication,
			},
			wantErr: true,
			setupMockR: func() {
				mockR.EXPECT().
					FetchGenre(fakeCtx, "anyname").
					Return(
						domain.GenreValidatable{},
						fakeErrorInternalApplication,
					)
			},
		},
		{
			name: "When name is not found",
			args: args{fakeCtx, fakeDoesNotExistName},
			want: returns{
				domain.GenreValidatable{},
				fmt.Errorf("%s: %w", fakeDoesNotExistName, logger.ErrNotFound),
			},
			wantErr: true,
			setupMockR: func() {
				mockR.EXPECT().
					FetchGenre(fakeCtx, strings.ToLower(strings.TrimSpace(fakeDoesNotExistName))).
					Return(
						domain.GenreValidatable{},
						sql.ErrNoRows,
					)
			},
		},
		{
			name: "When name is found",
			args: args{fakeCtx, fakeExistName},
			want: returns{
				domain.GenreValidatable{
					Name: fakeExistName,
				},
				nil,
			},
			wantErr: false,
			setupMockR: func() {
				mockR.EXPECT().
					FetchGenre(fakeCtx, fakeExistName).
					Return(
						domain.GenreValidatable{
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
			s := service.NewService(mockR)
			got, err := s.FetchGenre(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchGenre() error: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.genre) {
				t.Errorf("GetGenre() got: %v, want: %v", got, tt.want.genre)
			}
			if tt.wantErr && errors.Is(err, tt.want.err) {
				t.Errorf("GetGenre() got: %v, want: %v", err, tt.want.err)
			}
		})
	}
}
