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
)

func TestAddCastMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeName := strings.ToLower(faker.FirstName())
	type fields struct {
		r sqlboiler.Repository
	}
	type returns struct {
		c   models.CastMember
		err error
	}
	type args struct {
		dto crud.CastMemberDTO
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    returns
		wantErr bool
	}{
		{
			name:    "When CastMemberDTO is not provided",
			args:    args{crud.CastMemberDTO{}},
			want:    returns{models.CastMember{}, logger.ErrIsRequired},
			wantErr: true,
		},
		{
			name: "When the name in CastMemberDTO is blank",
			args: args{crud.CastMemberDTO{
				Name: "    ",
				Type: crud.Director,
			}},
			want:    returns{models.CastMember{}, logger.ErrIsRequired},
			wantErr: true,
		},
		{
			name: "When cast member type in CastMemberDTO is not validated",
			args: args{crud.CastMemberDTO{
				Name: faker.Name(),
				Type: -1,
			}},
			want:    returns{models.CastMember{}, logger.ErrIsNotValidated},
			wantErr: true,
		},
		{
			name: "When CastMemberDTO is right",
			args: args{crud.CastMemberDTO{
				Name: fakeName,
				Type: crud.Actor,
			}},
			want: returns{models.CastMember{
				Name: fakeName,
				Type: int16(crud.Actor),
			}, nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockR.EXPECT().
					AddCastMember(tt.args.dto).
					Return(tt.want.err)
			}
			s := crud.NewService(mockR)
			err := s.AddCastMember(tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCastMember() error = '%v', wantErr '%v'", err, tt.wantErr)
				return
			}
			if !errors.Is(err, tt.want.err) {
				t.Errorf("AddCastMember() got = '%v', want '%v'", err, tt.want.err)
			}
		})
	}
}

func Test_service_RemoveCastMember(t *testing.T) {
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
					RemoveCastMember(tt.args.name).
					Return(tt.want)
			} else if tt.name == "When name is found" {
				mockR.EXPECT().
					RemoveCastMember(tt.args.name).
					Return(tt.want)
			}
			s := crud.NewService(mockR)
			err := s.RemoveCastMember(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveCastMember() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("RemoveCastMember() got: '%v', want: '%v'", err, tt.want)
			}
		})
	}
}

func Test_service_UpdateCastMember(t *testing.T) {
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
		dto  crud.CastMemberDTO
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
				crud.CastMemberDTO{
					Name: faker.Name(),
					Type: crud.Actor,
				},
			},
			want:    logger.ErrIsRequired,
			wantErr: true,
		},
		{
			name:    "When CastMemberDTO is not provided",
			args:    args{fakeExistName, crud.CastMemberDTO{}},
			want:    logger.ErrIsRequired,
			wantErr: true,
		},
		{
			name: "When name is not found",
			args: args{
				fakeDoesNotExistName,
				crud.CastMemberDTO{
					Name: faker.FirstName(),
					Type: crud.Director,
				},
			},
			want:    fmt.Errorf("%s: %w", fakeDoesNotExistName, logger.ErrNotFound),
			wantErr: true,
		},
		{
			name: "When name is found and CastMemberDTO is provided",
			args: args{
				fakeExistName,
				crud.CastMemberDTO{
					Name: faker.FirstName(),
					Type: crud.Actor,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "When name is not found" || tt.name == "When name is found and CastMemberDTO is provided" {
				mockR.EXPECT().
					UpdateCastMember(tt.args.name, tt.args.dto).
					Return(tt.want)
			}
			s := crud.NewService(mockR)
			err := s.UpdateCastMember(tt.args.name, tt.args.dto)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCastMember() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !errors.Is(err, tt.want) {
				t.Errorf("UpdateCastMember() got: '%v', want: '%v'", err, tt.want)
			}
		})
	}
}

func Test_service_GetCastMembers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeCastMemberSlice := models.CastMemberSlice{
		&models.CastMember{
			Name: "Ana Silva",
		},
		&models.CastMember{
			Name: "João Batista",
		},
		&models.CastMember{
			Name: "Maria Alves",
		},
	}
	fakeLimit := len(fakeCastMemberSlice)
	type args struct {
		limit int
	}
	type returns struct {
		cs models.CastMemberSlice
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
			want:    returns{fakeCastMemberSlice, nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockR.EXPECT().
					GetCastMembers(tt.args.limit).
					Return(
						fakeCastMemberSlice,
						nil,
					)
			}
			s := crud.NewService(mockR)
			got, err := s.GetCastMembers(tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCastMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.cs) {
				t.Errorf("GetCastMembers() got = %v, want %v", got, tt.want.cs)
			}
			if !reflect.DeepEqual(err, tt.want.e) {
				t.Errorf("GetCastMembers() got = %v, want %v", err, tt.want.e)
			}
		})
	}
}

func Test_service_FetchCastMember(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockR := mock.NewMockRepository(ctrl)
	fakeDoesNotExistName := "fakeDoesNotExistName"
	fakeExistName := "João Batista"
	fakeErrorInternalApplication := fmt.Errorf("Service.FetchCastMember(): %w", logger.ErrInternalApplication)
	type args struct {
		name string
	}
	type returns struct {
		c models.CastMember
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
				models.CastMember{},
				fakeErrorInternalApplication,
			},
			wantErr: true,
			setupMockR: func() {
				mockR.EXPECT().
					FetchCastMember("anyName").
					Return(
						models.CastMember{},
						fakeErrorInternalApplication,
					)
			},
		},
		{
			name: "When name is not found",
			args: args{fakeDoesNotExistName},
			want: returns{
				models.CastMember{},
				fmt.Errorf("%s: %w", fakeDoesNotExistName, logger.ErrNotFound),
			},
			wantErr: true,
			setupMockR: func() {
				mockR.EXPECT().
					FetchCastMember(fakeDoesNotExistName).
					Return(
						models.CastMember{},
						sql.ErrNoRows,
					)
			},
		},
		{
			name: "When name is found",
			args: args{fakeExistName},
			want: returns{
				models.CastMember{
					Name: fakeExistName,
				},
				nil,
			},
			wantErr: false,
			setupMockR: func() {
				mockR.EXPECT().
					FetchCastMember(fakeExistName).
					Return(
						models.CastMember{
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
			got, err := s.FetchCastMember(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchCastMember() error: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want.c) {
				t.Errorf("GetCastMember() got: %v, want: %v", got, tt.want.c)
			}
			if tt.wantErr && errors.Is(err, tt.want.e) {
				t.Errorf("GetCastMember() got: %v, want: %v", err, tt.want.e)
			}
		})
	}
}
