package genre_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/selmison/code-micro-videos/pkg/genre/mock"

	"github.com/selmison/code-micro-videos/pkg/genre"
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
		ctx      context.Context
		newGenre genre.NewGenre
	}
	tests := []struct {
		name          string
		args          args
		expectedError error
		wantErr       bool
	}{
		{
			name:          "When Genre is not provided",
			args:          args{testdata.FakeCtx, genre.NewGenre{}},
			expectedError: fmt.Errorf("'%s' param %v", "newGenre", logger.ErrCouldNotBeEmpty),
			wantErr:       true,
		},
		{
			name: "When the Name in Genre is blank",
			args: args{testdata.FakeCtx, genre.NewGenre{
				Name: "    ",
			}},
			expectedError: fmt.Errorf("'Name' field %v", logger.ErrIsRequired),
			wantErr:       true,
		},
		{
			name: "When the Name in Genre already exists",
			args: args{
				testdata.FakeCtx,
				genre.NewGenre{
					Name:         testdata.FakeName,
					CategoriesId: []string{testdata.FakeExistentGenreId},
				},
			},
			expectedError: fmt.Errorf("name '%s' %v", testdata.FakeName, logger.ErrAlreadyExists),
			wantErr:       true,
		},
		{
			name: "When Genre is with wrong genres",
			args: args{testdata.FakeCtx,
				genre.NewGenre{
					Name:         testdata.FakeName,
					CategoriesId: []string{testdata.FakeNonExistentGenre.Id},
				},
			},
			expectedError: fmt.Errorf("none genre is %v", logger.ErrNotFound),
			wantErr:       true,
		},
		{
			name: "When Genre is right",
			args: args{
				testdata.FakeCtx,
				genre.NewGenre{
					Name:         testdata.FakeName,
					CategoriesId: []string{testdata.FakeExistentGenre.Id},
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
			fakeGenre := genre.Genre{
				Id:           fakeUUID,
				Name:         strings.ToLower(tt.args.newGenre.Name),
				CategoriesId: tt.args.newGenre.CategoriesId,
			}
			switch tt.name {
			case "When the Name in Genre already exists":
				mockIdGen.EXPECT().
					Generate().
					Return(fakeUUID, nil)
				mockRepo.EXPECT().
					Store(tt.args.ctx, fakeGenre).
					Return(tt.expectedError)
			case "When Genre is with wrong genres":
				mockIdGen.EXPECT().
					Generate().
					Return(fakeUUID, nil)
				mockRepo.EXPECT().
					Store(tt.args.ctx, fakeGenre).
					Return(tt.expectedError)
			case "When Genre is right":
				mockIdGen.EXPECT().
					Generate().
					Return(fakeUUID, nil)
				mockRepo.EXPECT().
					Store(tt.args.ctx, fakeGenre).
					Return(nil)
			}
			svc := genre.NewService(mockIdGen, mockRepo, mockLogger)
			got, err := svc.Create(tt.args.ctx, tt.args.newGenre)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.expectedError.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, fakeGenre, got, "they should be equal")
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
			svc := genre.NewService(mockIdGen, mockRepo, mockLogger)
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
		genres        []genre.Genre
		expectedError error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "List genres",
			args: args{testdata.FakeCtx},
			want: returns{
				genres: []genre.Genre{
					{
						Id:           testdata.FakeId,
						Name:         testdata.FakeName,
						CategoriesId: []string{testdata.FakeExistentGenre.Id},
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
				Return(tt.want.genres, nil)
			svc := genre.NewService(mockIdGen, mockRepo, mockLogger)
			got, err := svc.List(tt.args.ctx)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.want.expectedError.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.genres, got, "they should be equal")
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
		genre         genre.Genre
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
				genre.Genre{},
				fmt.Errorf("'%s' param %w", "id", logger.ErrCouldNotBeEmpty),
			},
			wantErr: true,
		},
		{
			name: "When id is not found",
			args: args{testdata.FakeCtx, testdata.FakeNonExistentCategoryId},
			want: returns{
				genre.Genre{},
				fakeErrorNonExistentCategoryId,
			},
			wantErr: true,
		},
		{
			name: "When id is found",
			args: args{testdata.FakeCtx, testdata.FakeNonExistentCategoryId},
			want: returns{
				testdata.FakeExistentGenre,
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
					Return(genre.Genre{}, fakeErrorNonExistentCategoryId)
			case "When id is found":
				mockRepo.EXPECT().
					GetOne(tt.args.ctx, tt.args.id).
					Return(tt.want.genre, nil)
			}
			svc := genre.NewService(mockIdGen, mockRepo, mockLogger)
			got, err := svc.Show(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.want.expectedError.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.genre, got, "they should be equal")
			}
		})
	}
}

func Test_Service_Update(t *testing.T) {
	type args struct {
		ctx         context.Context
		id          string
		updateGenre genre.UpdateGenre
	}
	fakeBlankName := "     "
	fakeErrNonExistentCategoryId := fmt.Errorf("%s: %w", testdata.FakeNonExistentCategoryId, logger.ErrNotFound)
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
				genre.UpdateGenre{
					Name: &testdata.FakeName,
				},
			},
			want:    fmt.Errorf("'id' param %w", logger.ErrCouldNotBeEmpty),
			wantErr: true,
		},
		{
			name: "When the Name in UpdateGenre is blank",
			args: args{
				testdata.FakeCtx,
				testdata.FakeId,
				genre.UpdateGenre{
					Name: &fakeBlankName,
				},
			},
			want:    fmt.Errorf("the %s field %w", "Name", logger.ErrCouldNotBeEmpty),
			wantErr: true,
		},
		{
			name: "When UpdateGenre is with wrong genres",
			args: args{
				testdata.FakeCtx,
				testdata.FakeId,
				genre.UpdateGenre{
					Name: &testdata.FakeName,
					CategoriesId: []string{
						testdata.FakeNonExistentCategoryId,
					},
				},
			},
			want:    fakeErrNonExistentCategoryId,
			wantErr: true,
		},
		{
			name: "When id is not found",
			args: args{testdata.FakeCtx,
				testdata.FakeId,
				genre.UpdateGenre{
					Name:         &testdata.FakeName,
					CategoriesId: []string{testdata.FakeExistentGenre.Id},
				},
			},
			want:    fakeErrorIdIsNotFound,
			wantErr: true,
		},
		{
			name: "When UpdateGenre is not provided",
			args: args{
				testdata.FakeCtx,
				testdata.FakeId,
				genre.UpdateGenre{}},
			want:    nil,
			wantErr: false,
		},
		{
			name: "When the Name in UpdateGenre is nil",
			args: args{
				testdata.FakeCtx,
				testdata.FakeId,
				genre.UpdateGenre{
					Name: nil,
				}},
			want:    nil,
			wantErr: false,
		},
		{
			name: "When everything is right",
			args: args{testdata.FakeCtx,
				testdata.FakeId,
				genre.UpdateGenre{
					Name:         &testdata.FakeName,
					CategoriesId: []string{testdata.FakeExistentGenre.Id},
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
			case "When UpdateGenre is with wrong genres":
				mockRepo.EXPECT().
					UpdateOne(tt.args.ctx, tt.args.id, tt.args.updateGenre).
					Return(fakeErrNonExistentCategoryId)
			case "When id is not found":
				mockRepo.EXPECT().
					UpdateOne(tt.args.ctx, tt.args.id, tt.args.updateGenre).
					Return(fakeErrorIdIsNotFound)
			case "When everything is right":
				mockRepo.EXPECT().
					UpdateOne(tt.args.ctx, tt.args.id, tt.args.updateGenre).
					Return(tt.want)
			}
			svc := genre.NewService(mockIdGen, mockRepo, mockLogger)
			err := svc.Update(tt.args.ctx, tt.args.id, tt.args.updateGenre)
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
