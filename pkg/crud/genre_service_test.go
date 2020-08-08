package crud_test

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/selmison/code-micro-videos/models"
	"github.com/selmison/code-micro-videos/pkg/crud"
	"github.com/selmison/code-micro-videos/pkg/crud/mock"
	"github.com/selmison/code-micro-videos/pkg/logger"
)

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
					FetchGenre("anyName").
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
					FetchGenre(fakeDoesNotExistName).
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
