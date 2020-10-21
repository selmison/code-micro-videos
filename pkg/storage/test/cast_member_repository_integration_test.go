package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/selmison/code-micro-videos/pkg/cast_member"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/pkg/storage/inmem"
	"github.com/selmison/code-micro-videos/testdata"
)

var fakeCtx = context.Background()

func Test_cast_member_Repository_Create(t *testing.T) {
	fakeExistCastMember, err := cast_member.NewCastMember(
		uuid.New().String(),
		cast_member.NewCastMemberDTO{
			Name: "fakeExistentName",
			Type: cast_member.Actor,
		})
	if err != nil {
		t.Fatalf("test: failed to generate CastMember: %v\n", err)
	}
	inmem.NewCastMemberRepository()
	teardownTestCase, repo, err := SetupTestCase(t, []cast_member.CastMember{fakeExistCastMember})
	if err != nil {
		t.Fatalf("test: failed to setup test case: %v\n", err)
	}
	defer teardownTestCase(t)
	fakeNonExistentCastMember, err := cast_member.NewCastMember(
		uuid.New().String(),
		cast_member.NewCastMemberDTO{
			Name: "fakeNonExistentName",
			Type: cast_member.Actor,
		})
	if err != nil {
		t.Fatalf("test: failed to generate CastMember: %v\n", err)
	}
	type args struct {
		castMember cast_member.CastMember
	}
	tests := []struct {
		name          string
		args          args
		expectedError error
		wantErr       bool
	}{
		{
			name:          "When id in CastMember already exists",
			args:          args{fakeExistCastMember},
			expectedError: fmt.Errorf("id '%s' %w", fakeExistCastMember.Id(), logger.ErrAlreadyExists),
			wantErr:       true,
		},
		{
			name:          "When CastMember is right",
			args:          args{fakeNonExistentCastMember},
			expectedError: nil,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Store(fakeCtx, tt.args.castMember)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.expectedError.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_cast_member_Repository_GetAll(t *testing.T) {
	teardownTestCase, repo, err := SetupTestCase(t, testdata.FakeCastMembers)
	if err != nil {
		t.Fatalf("test: failed to setup test case: %v\n", err)
	}
	defer teardownTestCase(t)
	type returns struct {
		castMembers []cast_member.CastMember
		err         error
	}
	tests := []struct {
		name           string
		expectedResult returns
		wantErr        bool
	}{
		{
			name:           "When everything is right",
			expectedResult: returns{testdata.FakeCastMembers, nil},
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		fakeCtx := context.Background()
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetAll(fakeCtx)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.expectedResult.err.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.expectedResult.castMembers, got, "they should be equal")
			}
		})
	}
}

func Test_cast_member_Repository_Get(t *testing.T) {
	teardownTestCase, repo, err := SetupTestCase(t, testdata.FakeCastMembers)
	if err != nil {
		t.Fatalf("test: failed to setup test case: %v\n", err)
	}
	defer teardownTestCase(t)
	const (
		fakeDoesNonExistentId = "fakeDoesNonExistentId"
	)
	fakeExistentCastMember := testdata.FakeCastMembers[0]
	fakeExistentCastMemberId := fakeExistentCastMember.Id()
	type args struct {
		id string
	}
	type returns struct {
		castMember cast_member.CastMember
		err        error
	}
	tests := []struct {
		name           string
		args           args
		expectedResult returns
		wantErr        bool
	}{
		{
			name: "When id is not found",
			args: args{fakeDoesNonExistentId},
			expectedResult: returns{
				nil,
				fmt.Errorf("%s: %w", fakeDoesNonExistentId, logger.ErrNotFound),
			},
			wantErr: true,
		},
		{
			name: "When id is found",
			args: args{fakeExistentCastMemberId},
			expectedResult: returns{
				fakeExistentCastMember,
				nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetOne(fakeCtx, tt.args.id)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.expectedResult.err.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.castMember, got, "they should be equal")
			}
		})
	}
}

func Test_cast_member_Repository_RemoveCastMember(t *testing.T) {
	teardownTestCase, repo, err := SetupTestCase(t, testdata.FakeCastMembers)
	if err != nil {
		t.Fatalf("test: failed to setup test case: %v\n", err)
	}
	defer teardownTestCase(t)
	fakeExistentId := testdata.FakeCastMembers[0].Id()
	const (
		fakeDoesNonExistentId = "fakeDoesNonExistentId"
	)
	type fields struct {
		ctx context.Context
	}
	type args struct {
		id string
	}
	tests := []struct {
		id          string
		fields      fields
		args        args
		expectedErr error
		wantErr     bool
	}{
		{
			id:          "When id is not found",
			args:        args{fakeDoesNonExistentId},
			expectedErr: fmt.Errorf("%s: %w", fakeDoesNonExistentId, logger.ErrNotFound),
			wantErr:     true,
		},
		{
			id:          "When id is found",
			args:        args{fakeExistentId},
			expectedErr: nil,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			err := repo.DeleteOne(fakeCtx, tt.args.id)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.expectedErr.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_cast_member_Repository_UpdateOne(t *testing.T) {
	teardownTestCase, repo, err := SetupTestCase(t, testdata.FakeCastMembers)
	if err != nil {
		t.Fatalf("test: failed to setup test case: %v\n", err)
	}
	defer teardownTestCase(t)
	fakeExistentId := testdata.FakeCastMembers[0].Id()
	var (
		fakeDoesNonExistentId    = "fakeDoesNonExistentId"
		fakeNewDoestNotExistName = "new_action"
	)
	type fields struct {
		ctx context.Context
	}
	type args struct {
		id               string
		updateCastMember cast_member.UpdateCastMemberDTO
	}
	tests := []struct {
		id          string
		fields      fields
		args        args
		expectedErr error
		wantErr     bool
	}{
		{
			id: "When id to update doesn't exist",
			args: args{
				fakeDoesNonExistentId,
				cast_member.UpdateCastMemberDTO{
					Name: &fakeNewDoestNotExistName,
				},
			},
			expectedErr: fmt.Errorf("%s: %w", fakeDoesNonExistentId, logger.ErrNotFound),
			wantErr:     true,
		},
		{
			id: "When id exists and UpdateCastMemberDTO is right",
			args: args{
				fakeExistentId,
				cast_member.UpdateCastMemberDTO{
					Name: &fakeNewDoestNotExistName,
				},
			},
			expectedErr: nil,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			err := repo.UpdateOne(fakeCtx, tt.args.id, tt.args.updateCastMember)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.expectedErr.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
