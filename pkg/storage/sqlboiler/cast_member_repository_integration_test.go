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

func TestRepository_AddCastMember(t *testing.T) {
	cfg, teardownTestCase, repository, err := setupTestCase(nil)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	const (
		fakeExistName        = "action"
		fakeDoesNotExistName = "fakeDoesNotExistName"
	)
	fakeExistCastMember := models.CastMember{
		ID:   uuid.New().String(),
		Name: fakeExistName,
	}
	type args struct {
		castMemberDTO crud.CastMemberDTO
	}
	type returns struct {
		castMember models.CastMember
		err        error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "When name in CastMemberDTO already exists",
			args: args{crud.CastMemberDTO{
				Name: fakeExistName,
			}},
			want: returns{
				models.CastMember{
					Name: fakeExistName,
				},
				nil,
			},
			wantErr: false,
		},
		{
			name: "When CastMemberDTO is right",
			args: args{
				crud.CastMemberDTO{
					Name: fakeDoesNotExistName,
				},
			},
			want: returns{
				models.CastMember{
					Name: fakeDoesNotExistName,
				},
				nil,
			},
			wantErr: false,
		},
	}
	db, err := sql.Open(cfg.DBDrive, cfg.DBConnStr)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("test: failed to close DB: %v\n", err)
		}
	}()
	ctx := context.Background()
	err = fakeExistCastMember.InsertG(ctx, boil.Infer())
	if err != nil {
		t.Errorf("test: insert castMember: %s", err)
		return
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.AddCastMember(tt.args.castMemberDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCastMember() error: %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(err, tt.want.err) {
				t.Errorf("AddCastMember() got: %v, want: %v", err, tt.want.err)
				return
			}
		})
	}
}

func TestRepository_GetCastMembers(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeCastMembers)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	maximum := len(testdata.FakeCastMembers)
	type args struct {
		limit int
	}
	type returns struct {
		castMembers models.CastMemberSlice
		e           error
		amount      int
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
			got, err := repository.GetCastMembers(tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCastMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want.amount {
				t.Errorf("GetCastMembers() len(got): %v, want: %d", len(got), tt.want.amount)
			}
		})
	}
}

func TestRepository_FetchCastMember(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeCastMembers)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	fakeExistName := testdata.FakeCastMembers[0].Name
	const (
		fakeDoesNotExistName = "fakeDoesNotExistName"
	)
	type args struct {
		name string
	}
	type returns struct {
		castMember models.CastMember
		e          error
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
				models.CastMember{},
				sql.ErrNoRows,
			},
			wantErr: true,
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repository.FetchCastMember(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchCastMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Name != tt.want.castMember.Name {
				t.Errorf("FetchCastMember() got: %q, want: %q", got.Name, tt.want.castMember.Name)
			}
		})
	}
}

func TestRepository_RemoveCastMember(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeCastMembers)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	fakeExistName := testdata.FakeCastMembers[0].Name
	const (
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
			err := repository.RemoveCastMember(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveCastMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != tt.want {
				t.Errorf("RemoveCastMember() got: %s, want: %q", err, tt.want)
			}
		})
	}
}

func TestRepository_UpdateCastMember(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(testdata.FakeCastMembers)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	fakeExistName := testdata.FakeCastMembers[0].Name
	const (
		fakeDoesNotExistName     = "fakeDoesNotExistName"
		fakeNewDoestNotExistName = "new_action"
	)
	type fields struct {
		ctx context.Context
	}
	type args struct {
		name          string
		castMemberDTO crud.CastMemberDTO
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
				crud.CastMemberDTO{
					Name: fakeNewDoestNotExistName,
				},
			},
			want:    sql.ErrNoRows,
			wantErr: true,
		},
		{
			name: "When name exists and CastMemberDTO is right",
			args: args{
				fakeExistName,
				crud.CastMemberDTO{
					Name: fakeNewDoestNotExistName,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repository.UpdateCastMember(tt.args.name, tt.args.castMemberDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateCastMember() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(err, tt.want) {
				t.Errorf("UpdateCastMember() got: %v, want: %v", err, tt.want)
			}
		})
	}
}

func TestCastMember_isValidUUIDHook(t *testing.T) {
	_, teardownTestCase, repository, err := setupTestCase(nil)
	if err != nil {
		t.Errorf("test: failed to setup test case: %v\n", err)
		return
	}
	defer teardownTestCase(t)
	type args struct {
		castMember models.CastMember
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
				models.CastMember{
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
				models.CastMember{
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
			err := tt.args.castMember.InsertG(repository.ctx, boil.Infer())
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
