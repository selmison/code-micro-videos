package cast_member_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/selmison/code-micro-videos/pkg/cast_member/mock"

	"github.com/selmison/code-micro-videos/pkg/cast_member"
	mockId "github.com/selmison/code-micro-videos/pkg/id_generator/mock"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/testdata"
)

type mockLogger struct{}

func (l mockLogger) Log(_ ...interface{}) error {
	return nil
}

func Test_Service_Create(t *testing.T) {
	type args struct {
		ctx           context.Context
		newCastMember cast_member.NewCastMemberDTO
	}
	fakeErrNonValidatedCastMemberType := fmt.Errorf("cast member type %v", logger.ErrIsNotValidated)
	tests := []struct {
		name          string
		args          args
		expectedError error
		wantErr       bool
	}{
		{
			name:          "When CastMember is not provided",
			args:          args{testdata.FakeCtx, cast_member.NewCastMemberDTO{}},
			expectedError: fmt.Errorf("'%s' param %v", "newCastMember", logger.ErrCouldNotBeEmpty),
			wantErr:       true,
		},
		{
			name: "When the Name in CastMember is blank",
			args: args{testdata.FakeCtx, cast_member.NewCastMemberDTO{
				Name: "    ",
			}},
			expectedError: fmt.Errorf("'Name' field %v", logger.ErrCouldNotBeEmpty),
			wantErr:       true,
		},
		{
			name: "When the Name in CastMember already exists",
			args: args{
				testdata.FakeCtx,
				cast_member.NewCastMemberDTO{
					Name: testdata.FakeName,
					Type: cast_member.Actor,
				},
			},
			expectedError: fmt.Errorf("name '%s' %v", testdata.FakeName, logger.ErrAlreadyExists),
			wantErr:       true,
		},
		{
			name: "When CastMember is with wrong CastMemberType",
			args: args{testdata.FakeCtx,
				cast_member.NewCastMemberDTO{
					Name: testdata.FakeName,
					Type: 111,
				},
			},
			expectedError: fakeErrNonValidatedCastMemberType,
			wantErr:       true,
		},
		{
			name: "When CastMember is right",
			args: args{
				testdata.FakeCtx,
				cast_member.NewCastMemberDTO{
					Name: testdata.FakeName,
					Type: cast_member.Actor,
				},
			},
			wantErr: false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock.NewMockRepository(ctrl)
	mockIdGen := mockId.NewMockIdGenerator(ctrl)
	mockLogger := mockLogger{}
	fakeUUID := uuid.New().String()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeCastMember, _ := cast_member.NewCastMember(
				fakeUUID,
				cast_member.NewCastMemberDTO{
					Name: tt.args.newCastMember.Name,
					Type: tt.args.newCastMember.Type,
				})
			switch tt.name {
			case "When the Name in CastMember already exists":
				mockIdGen.EXPECT().
					Generate().
					Return(fakeUUID, nil)
				mockRepo.EXPECT().
					Store(tt.args.ctx, fakeCastMember).
					Return(tt.expectedError)
			case "When CastMember is right":
				mockIdGen.EXPECT().
					Generate().
					Return(fakeUUID, nil)
				mockRepo.EXPECT().
					Store(tt.args.ctx, fakeCastMember).
					Return(nil)
			default:
				mockIdGen.EXPECT().
					Generate().
					Return(fakeUUID, nil)
			}
			svc := cast_member.NewService(mockIdGen, mockRepo, mockLogger)
			got, err := svc.Create(tt.args.ctx, tt.args.newCastMember)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.expectedError.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, fakeCastMember, got, "they should be equal")
			}
		})
	}
}

func Test_Service_Destroy(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name          string
		args          args
		expectedError error
		wantErr       bool
	}{
		{
			name:          "When id is blank",
			args:          args{testdata.FakeCtx, "     "},
			expectedError: fmt.Errorf("'%s' param %v", "id", logger.ErrCouldNotBeEmpty),
			wantErr:       true,
		},
		{
			name:          "When id is not found",
			args:          args{testdata.FakeCtx, testdata.FakeId},
			expectedError: fmt.Errorf("%s: %v", testdata.FakeId, logger.ErrNotFound),
			wantErr:       true,
		},
		{
			name:    "When id is found",
			args:    args{testdata.FakeCtx, testdata.FakeId},
			wantErr: false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock.NewMockRepository(ctrl)
	mockIdGen := mockId.NewMockIdGenerator(ctrl)
	mockLogger := mockLogger{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "When id is not found":
				mockRepo.EXPECT().
					DeleteOne(tt.args.ctx, tt.args.id).
					Return(tt.expectedError)
			case "When id is found":
				mockRepo.EXPECT().
					DeleteOne(tt.args.ctx, tt.args.id).
					Return(nil)
			}
			svc := cast_member.NewService(mockIdGen, mockRepo, mockLogger)
			err := svc.Destroy(tt.args.ctx, tt.args.id)
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

func Test_Service_List(t *testing.T) {
	fakeCastMember, _ := cast_member.NewCastMember(
		testdata.FakeId,
		cast_member.NewCastMemberDTO{
			Name: testdata.FakeName,
			Type: cast_member.Actor,
		})
	type args struct {
		ctx context.Context
	}
	type returns struct {
		castMembers   []cast_member.CastMember
		expectedError error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "List castMembers",
			args: args{testdata.FakeCtx},
			want: returns{
				castMembers: []cast_member.CastMember{
					fakeCastMember,
				},
			},
			wantErr: false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock.NewMockRepository(ctrl)
	mockIdGen := mockId.NewMockIdGenerator(ctrl)
	mockLogger := mockLogger{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().
				GetAll(tt.args.ctx).
				Return(tt.want.castMembers, nil)
			svc := cast_member.NewService(mockIdGen, mockRepo, mockLogger)
			got, err := svc.List(tt.args.ctx)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.want.expectedError.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.castMembers, got, "they should be equal")
			}
		})
	}
}

func Test_Service_Show(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	type returns struct {
		castMember    cast_member.CastMember
		expectedError error
	}
	fakeErrorNonExistentCategoryId := fmt.Errorf("%s: %w", testdata.FakeNonExistentCategoryId, logger.ErrNotFound)
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "When id is blank",
			args: args{
				testdata.FakeCtx,
				"     ",
			},
			want: returns{
				nil,
				fmt.Errorf("'%s' param %w", "id", logger.ErrCouldNotBeEmpty),
			},
			wantErr: true,
		},
		{
			name: "When id is not found",
			args: args{testdata.FakeCtx, testdata.FakeNonExistentCategoryId},
			want: returns{
				nil,
				fakeErrorNonExistentCategoryId,
			},
			wantErr: true,
		},
		{
			name: "When id is found",
			args: args{testdata.FakeCtx, testdata.FakeNonExistentCategoryId},
			want: returns{
				testdata.FakeExistentCastMember,
				nil,
			},
			wantErr: false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock.NewMockRepository(ctrl)
	mockIdGen := mockId.NewMockIdGenerator(ctrl)
	mockLogger := mockLogger{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "When id is not found":
				mockRepo.EXPECT().
					GetOne(tt.args.ctx, tt.args.id).
					Return(nil, fakeErrorNonExistentCategoryId)
			case "When id is found":
				mockRepo.EXPECT().
					GetOne(tt.args.ctx, tt.args.id).
					Return(tt.want.castMember, nil)
			}
			svc := cast_member.NewService(mockIdGen, mockRepo, mockLogger)
			got, err := svc.Show(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.want.expectedError.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.castMember, got, "they should be equal")
			}
		})
	}
}

func Test_Service_Update(t *testing.T) {
	type args struct {
		ctx              context.Context
		id               string
		updateCastMember cast_member.UpdateCastMemberDTO
	}
	fakeBlankName := "     "
	fakeValidatedType := cast_member.Actor
	fakeNonValidatedType := cast_member.CastMemberType(111)
	fakeErrNonValidatedCastMemberType := fmt.Errorf("cast member type %v", logger.ErrIsNotValidated)
	fakeErrorIdIsNotFound := fmt.Errorf("%s: %w", testdata.FakeId, logger.ErrNotFound)
	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "When id is blank",
			args: args{
				testdata.FakeCtx,
				"     ",
				cast_member.UpdateCastMemberDTO{
					Name: &testdata.FakeName,
				},
			},
			want:    fmt.Errorf("'id' param %w", logger.ErrCouldNotBeEmpty),
			wantErr: true,
		},
		{
			name: "When the Name in UpdateCastMemberDTO is blank",
			args: args{
				testdata.FakeCtx,
				testdata.FakeId,
				cast_member.UpdateCastMemberDTO{
					Name: &fakeBlankName,
				},
			},
			want:    fmt.Errorf("the %s field %w", "Name", logger.ErrCouldNotBeEmpty),
			wantErr: true,
		},
		{
			name: "When UpdateCastMemberDTO is with wrong CastMemberType",
			args: args{
				testdata.FakeCtx,
				testdata.FakeId,
				cast_member.UpdateCastMemberDTO{
					Name: &testdata.FakeName,
					Type: &fakeNonValidatedType,
				},
			},
			want:    fakeErrNonValidatedCastMemberType,
			wantErr: true,
		},
		{
			name: "When id is not found",
			args: args{testdata.FakeCtx,
				testdata.FakeId,
				cast_member.UpdateCastMemberDTO{
					Name: &testdata.FakeName,
					Type: &fakeValidatedType,
				},
			},
			want:    fakeErrorIdIsNotFound,
			wantErr: true,
		},
		{
			name: "When UpdateCastMemberDTO is not provided",
			args: args{
				testdata.FakeCtx,
				testdata.FakeId,
				cast_member.UpdateCastMemberDTO{}},
			want:    nil,
			wantErr: false,
		},
		{
			name: "When the Name in UpdateCastMemberDTO is nil",
			args: args{
				testdata.FakeCtx,
				testdata.FakeId,
				cast_member.UpdateCastMemberDTO{
					Name: nil,
				}},
			want:    nil,
			wantErr: false,
		},
		{
			name: "When everything is right",
			args: args{testdata.FakeCtx,
				testdata.FakeId,
				cast_member.UpdateCastMemberDTO{
					Name: &testdata.FakeName,
					Type: &fakeValidatedType,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock.NewMockRepository(ctrl)
	mockIdGen := mockId.NewMockIdGenerator(ctrl)
	mockLogger := mockLogger{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "When UpdateCastMemberDTO is with wrong castMembers":
				mockRepo.EXPECT().
					UpdateOne(tt.args.ctx, tt.args.id, tt.args.updateCastMember).
					Return(fakeErrNonValidatedCastMemberType)
			case "When id is not found":
				mockRepo.EXPECT().
					UpdateOne(tt.args.ctx, tt.args.id, tt.args.updateCastMember).
					Return(fakeErrorIdIsNotFound)
			case "When everything is right":
				mockRepo.EXPECT().
					UpdateOne(tt.args.ctx, tt.args.id, tt.args.updateCastMember).
					Return(tt.want)
			}
			svc := cast_member.NewService(mockIdGen, mockRepo, mockLogger)
			err := svc.Update(tt.args.ctx, tt.args.id, tt.args.updateCastMember)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.want.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
