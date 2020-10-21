package category_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/selmison/code-micro-videos/pkg/category"
	"github.com/selmison/code-micro-videos/pkg/category/mock"
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
		ctx         context.Context
		newCategory category.NewCategory
	}
	tests := []struct {
		name          string
		args          args
		expectedError error
		wantErr       bool
	}{
		{
			name:          "When Category is not provided",
			args:          args{testdata.FakeCtx, category.NewCategory{}},
			expectedError: fmt.Errorf("'%s' param %v", "newCategory", logger.ErrCouldNotBeEmpty),
			wantErr:       true,
		},
		{
			name: "When the Name in Category is blank",
			args: args{
				testdata.FakeCtx,
				category.NewCategory{
					Name:        "    ",
					Description: testdata.FakeDesc,
				},
			},
			expectedError: fmt.Errorf("'Name' field %v", logger.ErrIsRequired),
			wantErr:       true,
		},
		{
			name: "When the Name in Category already exists",
			args: args{testdata.FakeCtx, category.NewCategory{
				Name:        testdata.FakeName,
				Description: testdata.FakeDesc,
				GenresId:    []string{testdata.FakeExistentGenre.Id},
			},
			},
			expectedError: fmt.Errorf("name '%s' %v", testdata.FakeName, logger.ErrAlreadyExists),
			wantErr:       true,
		},
		{
			name: "When Category is with wrong genres",
			args: args{testdata.FakeCtx,
				category.NewCategory{
					Name:        testdata.FakeName,
					Description: testdata.FakeDesc,
					GenresId:    []string{testdata.FakeNonExistentGenre.Id},
				},
			},
			expectedError: fmt.Errorf("none genre is %v", logger.ErrNotFound),
			wantErr:       true,
		},
		{
			name: "When Category is right",
			args: args{
				testdata.FakeCtx,
				category.NewCategory{
					Name:        testdata.FakeName,
					Description: testdata.FakeDesc,
					GenresId:    []string{testdata.FakeExistentGenre.Id},
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
			fakeUUID := uuid.New().String()
			fakeCategory := category.Category{
				Id:          fakeUUID,
				Name:        strings.ToLower(tt.args.newCategory.Name),
				Description: tt.args.newCategory.Description,
				GenresId:    tt.args.newCategory.GenresId,
			}
			switch tt.name {
			case "When the Name in Category already exists":
				mockIdGen.EXPECT().
					Generate().
					Return(fakeUUID, nil)
				mockRepo.EXPECT().
					Store(tt.args.ctx, fakeCategory).
					Return(tt.expectedError)
			case "When Category is with wrong genres":
				mockIdGen.EXPECT().
					Generate().
					Return(fakeUUID, nil)
				mockRepo.EXPECT().
					Store(tt.args.ctx, fakeCategory).
					Return(tt.expectedError)
			case "When Category is right":
				mockIdGen.EXPECT().
					Generate().
					Return(fakeUUID, nil)
				mockRepo.EXPECT().
					Store(tt.args.ctx, fakeCategory).
					Return(nil)
			}
			svc := category.NewService(mockIdGen, mockRepo, mockLogger)
			got, err := svc.Create(tt.args.ctx, tt.args.newCategory)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.expectedError.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, fakeCategory, got, "they should be equal")
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
			expectedError: fmt.Errorf("'%s' params %v", "id", logger.ErrCouldNotBeEmpty),
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
			svc := category.NewService(mockIdGen, mockRepo, mockLogger)
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
	type args struct {
		ctx context.Context
	}
	type returns struct {
		categories    []category.Category
		expectedError error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "List categories",
			args: args{testdata.FakeCtx},
			want: returns{
				categories: []category.Category{
					{
						Id:          testdata.FakeId,
						Name:        testdata.FakeName,
						Description: testdata.FakeDesc,
						GenresId:    []string{testdata.FakeExistentGenre.Id},
					},
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
				Return(tt.want.categories, nil)
			svc := category.NewService(mockIdGen, mockRepo, mockLogger)
			got, err := svc.List(tt.args.ctx)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.want.expectedError.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.categories, got, "they should be equal")
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
		category      category.Category
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
				category.Category{},
				fmt.Errorf("'%s' param %w", "id", logger.ErrCouldNotBeEmpty),
			},
			wantErr: true,
		},
		{
			name: "When id is not found",
			args: args{testdata.FakeCtx, testdata.FakeNonExistentCategoryId},
			want: returns{
				category.Category{},
				fakeErrorNonExistentCategoryId,
			},
			wantErr: true,
		},
		{
			name: "When id is found",
			args: args{testdata.FakeCtx, testdata.FakeExistentCategoryId},
			want: returns{
				testdata.FakeExistentCategory,
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
					Return(category.Category{}, fakeErrorNonExistentCategoryId)
			case "When id is found":
				mockRepo.EXPECT().
					GetOne(tt.args.ctx, tt.args.id).
					Return(tt.want.category, nil)
			}
			svc := category.NewService(mockIdGen, mockRepo, mockLogger)
			got, err := svc.Show(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.want.expectedError.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.category, got, "they should be equal")
			}
		})
	}
}

func Test_Service_Update(t *testing.T) {
	type args struct {
		ctx            context.Context
		id             string
		updateCategory category.UpdateCategory
	}
	fakeBlankName := "     "
	fakeErrorNonExistentCategoryId := fmt.Errorf("%s: %w", testdata.FakeNonExistentCategoryId, logger.ErrNotFound)
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
				category.UpdateCategory{
					Name:        &testdata.FakeName,
					Description: &testdata.FakeDesc,
				},
			},
			want:    fmt.Errorf("'id' params %w", logger.ErrCouldNotBeEmpty),
			wantErr: true,
		},
		{
			name: "When the Name in UpdateCategory is blank",
			args: args{
				testdata.FakeCtx,
				testdata.FakeId,
				category.UpdateCategory{
					Name:        &fakeBlankName,
					Description: &testdata.FakeDesc,
				},
			},
			want:    fmt.Errorf("the %s field %w", "Name", logger.ErrCouldNotBeEmpty),
			wantErr: true,
		},
		{
			name: "When UpdateCategory is with wrong genres",
			args: args{
				testdata.FakeCtx,
				testdata.FakeId,
				category.UpdateCategory{
					Name:        &testdata.FakeName,
					Description: &testdata.FakeDesc,
					GenresId: []string{
						testdata.FakeNonExistentCategoryId,
					},
				},
			},
			want:    fakeErrorNonExistentCategoryId,
			wantErr: true,
		},
		{
			name: "When id is not found",
			args: args{testdata.FakeCtx,
				testdata.FakeId,
				category.UpdateCategory{
					Name:        &testdata.FakeName,
					Description: &testdata.FakeDesc,
					GenresId:    []string{testdata.FakeExistentGenre.Id},
				},
			},
			want:    fakeErrorIdIsNotFound,
			wantErr: true,
		},
		{
			name: "When UpdateCategory is not provided",
			args: args{
				testdata.FakeCtx,
				testdata.FakeId,
				category.UpdateCategory{}},
			want:    nil,
			wantErr: false,
		},
		{
			name: "When the Name in UpdateCategory is nil",
			args: args{
				testdata.FakeCtx,
				testdata.FakeId,
				category.UpdateCategory{
					Name:        nil,
					Description: &testdata.FakeDesc,
				}},
			want:    nil,
			wantErr: false,
		},
		{
			name: "When everything is right",
			args: args{testdata.FakeCtx,
				testdata.FakeId,
				category.UpdateCategory{
					Name:        &testdata.FakeName,
					Description: &testdata.FakeDesc,
					GenresId:    []string{testdata.FakeExistentGenre.Id},
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
			case "When UpdateCategory is with wrong genres":
				mockRepo.EXPECT().
					UpdateOne(tt.args.ctx, tt.args.id, tt.args.updateCategory).
					Return(fakeErrorNonExistentCategoryId)
			case "When id is not found":
				mockRepo.EXPECT().
					UpdateOne(tt.args.ctx, tt.args.id, tt.args.updateCategory).
					Return(fakeErrorIdIsNotFound)
			case "When the Name in UpdateCategory is nil",
				"When everything is right":
				mockRepo.EXPECT().
					UpdateOne(tt.args.ctx, tt.args.id, tt.args.updateCategory).
					Return(tt.want)
			}
			svc := category.NewService(mockIdGen, mockRepo, mockLogger)
			err := svc.Update(tt.args.ctx, tt.args.id, tt.args.updateCategory)
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
