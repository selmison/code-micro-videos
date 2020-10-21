package video_test

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	mockId "github.com/selmison/code-micro-videos/pkg/id_generator/mock"
	"github.com/selmison/code-micro-videos/pkg/logger"
	"github.com/selmison/code-micro-videos/pkg/video"
	"github.com/selmison/code-micro-videos/pkg/video/mock"
	"github.com/selmison/code-micro-videos/testdata"
)

var (
	fakeTitle                    = faker.Name()
	fakeYearLaunched       int16 = 2020
	fakeDuration           int16 = 90
	fakeRating                   = video.TwelveRating
	fakeNotValidatedRating       = video.VideoRating(111)
)

type mockLogger struct{}

func (l mockLogger) Log(_ ...interface{}) error {
	return nil
}

func TestVideoRating_String(t *testing.T) {
	tests := []struct {
		name string
		v    video.VideoRating
		want string
	}{
		{
			name: "When VideoRating is FreeRating",
			v:    video.FreeRating,
			want: "Free",
		},
		{
			name: "When VideoRating is TenRating",
			v:    video.TenRating,
			want: "10",
		},
		{
			name: "When VideoRating is TwelveRating",
			v:    video.TwelveRating,
			want: "12",
		},
		{
			name: "When VideoRating is FourteenRating",
			v:    video.FourteenRating,
			want: "14",
		},
		{
			name: "When VideoRating is SixteenRating",
			v:    video.SixteenRating,
			want: "16",
		},
		{
			name: "When VideoRating is EighteenRating",
			v:    video.EighteenRating,
			want: "18",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v.String()
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestVideoRating_Validate(t *testing.T) {
	tests := []struct {
		name          string
		v             video.VideoRating
		expectedError string
		wantErr       bool
	}{
		{
			name:          "When VideoRating is not validated",
			v:             fakeNotValidatedRating,
			expectedError: fmt.Sprintf("video rating %v", logger.ErrIsNotValidated),
			wantErr:       true,
		},
		{
			name:    "When VideoRating is FreeRating",
			v:       video.FreeRating,
			wantErr: false,
		},
		{
			name:    "When VideoRating is TenRating",
			v:       video.TenRating,
			wantErr: false,
		},
		{
			name:    "When VideoRating is TwelveRating",
			v:       video.TwelveRating,
			wantErr: false,
		},
		{
			name:    "When VideoRating is FourteenRating",
			v:       video.FourteenRating,
			wantErr: false,
		},
		{
			name:    "When VideoRating is SixteenRating",
			v:       video.SixteenRating,
			wantErr: false,
		},
		{
			name:    "When VideoRating is EighteenRating",
			v:       video.EighteenRating,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.v.Validate()
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.expectedError, "they should be equal")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewVideo_Validate(t *testing.T) {
	type fields struct {
		Title            string
		Description      string
		YearLaunched     *int16
		Opened           bool
		Rating           *video.VideoRating
		Duration         *int16
		CategoriesId     []string
		GenresId         []string
		VideoFileHandler *multipart.FileHeader
	}
	tests := []struct {
		name          string
		fields        fields
		expectedError string
		wantErr       bool
	}{
		{
			name: "When the Title is blank",
			fields: fields{
				Title:        "    ",
				YearLaunched: &fakeYearLaunched,
				Opened:       false,
				Rating:       &fakeRating,
				Duration:     &fakeDuration,
			},
			expectedError: fmt.Sprintf("'Title' field %v", logger.ErrIsRequired),
			wantErr:       true,
		},
		{
			name: "When the YearLaunched is blank",
			fields: fields{
				Title:        fakeTitle,
				YearLaunched: nil,
				Opened:       false,
				Rating:       &fakeRating,
				Duration:     &fakeDuration,
			},
			expectedError: fmt.Sprintf("'YearLaunched' field %v", logger.ErrIsRequired),
			wantErr:       true,
		},
		{
			name: "When the Rating is blank",
			fields: fields{
				Title:        fakeTitle,
				YearLaunched: &fakeYearLaunched,
				Opened:       false,
				Rating:       nil,
				Duration:     &fakeDuration,
			},
			expectedError: fmt.Sprintf("'Rating' field %v", logger.ErrIsRequired),
			wantErr:       true,
		},
		{
			name: "When the Duration is blank",
			fields: fields{
				Title:        fakeTitle,
				YearLaunched: &fakeYearLaunched,
				Opened:       false,
				Rating:       &fakeRating,
				Duration:     nil,
			},
			expectedError: fmt.Sprintf("'Duration' field %v", logger.ErrIsRequired),
			wantErr:       true,
		},
		{
			name: "When Video is without genres and categories",
			fields: fields{
				Title:        fakeTitle,
				Description:  testdata.FakeDesc,
				YearLaunched: &fakeYearLaunched,
				Opened:       false,
				Rating:       &fakeRating,
				Duration:     &fakeDuration,
			},
			expectedError: fmt.Sprintf("'CategoriesId' field %v", logger.ErrIsRequired),
			wantErr:       true,
		},
		{
			name: "When VideoRating is not validated",
			fields: fields{
				Title:        fakeTitle,
				Description:  testdata.FakeDesc,
				YearLaunched: &fakeYearLaunched,
				Opened:       false,
				Rating:       &fakeNotValidatedRating,
				Duration:     &fakeDuration,
				GenresId: []string{
					faker.UUIDHyphenated(),
				},
				CategoriesId: []string{
					faker.UUIDHyphenated(),
				},
			},
			expectedError: fmt.Sprintf("video rating %v", logger.ErrIsNotValidated),
			wantErr:       true,
		},
		{
			name: "When everything is right",
			fields: fields{
				Title:        fakeTitle,
				Description:  testdata.FakeDesc,
				YearLaunched: &fakeYearLaunched,
				Opened:       false,
				Rating:       &fakeRating,
				Duration:     &fakeDuration,
				GenresId: []string{
					faker.UUIDHyphenated(),
				},
				CategoriesId: []string{
					faker.UUIDHyphenated(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &video.NewVideo{
				Title:            tt.fields.Title,
				Description:      tt.fields.Description,
				YearLaunched:     tt.fields.YearLaunched,
				Opened:           tt.fields.Opened,
				Rating:           tt.fields.Rating,
				Duration:         tt.fields.Duration,
				CategoriesId:     tt.fields.CategoriesId,
				GenresId:         tt.fields.GenresId,
				VideoFileHandler: tt.fields.VideoFileHandler,
			}
			err := v.Validate()
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.expectedError, "they should be equal")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_Service_Create(t *testing.T) {
	type returns struct {
		video video.Video
		id    string
		err   error
	}
	type args struct {
		ctx      context.Context
		newVideo video.NewVideo
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "When NewVideo is not provided",
			args: args{
				testdata.FakeCtx,
				video.NewVideo{},
			},
			want:    returns{err: fmt.Errorf("'newVideo' param %w", logger.ErrCouldNotBeEmpty)},
			wantErr: true,
		},
		{
			name: "When the Title in NewVideo is blank",
			args: args{
				testdata.FakeCtx,
				video.NewVideo{
					Title:        "    ",
					YearLaunched: &fakeYearLaunched,
					Opened:       false,
					Rating:       &fakeRating,
					Duration:     &fakeDuration,
				},
			},
			want:    returns{err: fmt.Errorf("'Title' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When the YearLaunched in NewVideo is blank",
			args: args{
				testdata.FakeCtx,
				video.NewVideo{
					Title:        fakeTitle,
					YearLaunched: nil,
					Opened:       false,
					Rating:       &fakeRating,
					Duration:     &fakeDuration,
				},
			},
			want:    returns{err: fmt.Errorf("'YearLaunched' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When the Rating in NewVideo is blank",
			args: args{
				testdata.FakeCtx,
				video.NewVideo{
					Title:        fakeTitle,
					YearLaunched: &fakeYearLaunched,
					Opened:       false,
					Rating:       nil,
					Duration:     &fakeDuration,
				},
			},
			want:    returns{err: fmt.Errorf("'Rating' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When the Duration in NewVideo is blank",
			args: args{
				testdata.FakeCtx,
				video.NewVideo{
					Title:        fakeTitle,
					YearLaunched: &fakeYearLaunched,
					Opened:       false,
					Rating:       &fakeRating,
					Duration:     nil,
				},
			},
			want:    returns{err: fmt.Errorf("'Duration' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When NewVideo is with wrong genres and genres",
			args: args{
				testdata.FakeCtx,
				video.NewVideo{
					Title:        fakeTitle,
					Description:  testdata.FakeDesc,
					YearLaunched: &fakeYearLaunched,
					Opened:       false,
					Rating:       &fakeRating,
					Duration:     &fakeDuration,
					GenresId:     []string{testdata.FakeExistentCategoryId},
					CategoriesId: []string{testdata.FakeNonExistentCategoryId},
				},
			},
			want:    returns{err: logger.ErrIsRequired},
			wantErr: true,
		},
		{
			name: "When NewVideo is without genres and genres",
			args: args{
				testdata.FakeCtx,
				video.NewVideo{
					Title:        fakeTitle,
					Description:  testdata.FakeDesc,
					YearLaunched: &fakeYearLaunched,
					Opened:       false,
					Rating:       &fakeRating,
					Duration:     &fakeDuration,
				},
			},
			want:    returns{err: fmt.Errorf("'CategoriesId' field %w", logger.ErrIsRequired)},
			wantErr: true,
		},
		{
			name: "When NewVideo is right",
			args: args{
				testdata.FakeCtx,
				video.NewVideo{
					Title:        fakeTitle,
					Description:  testdata.FakeDesc,
					YearLaunched: &fakeYearLaunched,
					Opened:       false,
					Rating:       &fakeRating,
					Duration:     &fakeDuration,
					GenresId:     []string{testdata.FakeExistentGenreId},
					CategoriesId: []string{testdata.FakeExistentCategoryId},
				}},
			want: returns{video: video.Video{
				Title:        fakeTitle,
				Description:  testdata.FakeDesc,
				Opened:       false,
				YearLaunched: fakeYearLaunched,
				Rating:       fakeRating,
				Duration:     fakeDuration,
			},
				err: nil},
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
			fakeVideo := tt.args.newVideo.ToVideo(fakeUUID)
			if !tt.wantErr || tt.name == "When NewVideo is with wrong genres and genres" {
				mockIdGen.EXPECT().
					Generate().
					Return(fakeUUID, nil)
				mockRepo.EXPECT().
					Store(tt.args.ctx, fakeVideo).
					Return(tt.want.err)
			}
			s := video.NewService(mockIdGen, mockRepo, mockLogger)
			got, err := s.Create(tt.args.ctx, tt.args.newVideo)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.want.err.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, fakeVideo, got, "they should be equal")
			}
		})
	}
}

func Test_Service_Destroy(t *testing.T) {
	const fakeDoesNotExistTitle = "fakeDoesNotExistTitle"
	type args struct {
		ctx context.Context
		id  string
	}
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
			},
			want:    fmt.Errorf("'id' %w", logger.ErrCouldNotBeEmpty),
			wantErr: true,
		},
		{
			name: "When id is not found",
			args: args{
				testdata.FakeCtx,
				testdata.FakeId,
			},
			want:    fmt.Errorf("%s: %w", fakeDoesNotExistTitle, logger.ErrNotFound),
			wantErr: true,
		},
		{
			name: "When id is found",
			args: args{
				testdata.FakeCtx,
				testdata.FakeId,
			},
			want:    nil,
			wantErr: false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockIdGen := mockId.NewMockIdGenerator(ctrl)
	mockRepo := mock.NewMockRepository(ctrl)
	mockLogger := mockLogger{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "When id is not found" {
				tt.args.id = strings.ToLower(tt.args.id)
				mockRepo.EXPECT().
					DeleteOne(tt.args.ctx, tt.args.id).
					Return(tt.want)
			} else if tt.name == "When id is found" {
				mockRepo.EXPECT().
					DeleteOne(tt.args.ctx, tt.args.id).
					Return(tt.want)
			}
			s := video.NewService(mockIdGen, mockRepo, mockLogger)
			err := s.Destroy(tt.args.ctx, tt.args.id)
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

func Test_Service_List(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type returns struct {
		videos []video.Video
		err    error
	}
	tests := []struct {
		name    string
		args    args
		want    returns
		wantErr bool
	}{
		{
			name: "When everything is right",
			args: args{testdata.FakeCtx},
			want: returns{[]video.Video{
				{
					Id:           faker.UUIDHyphenated(),
					Title:        faker.Sentence(),
					Description:  faker.Sentence(),
					YearLaunched: fakeYearLaunched,
					Opened:       false,
					Rating:       fakeRating,
					Duration:     fakeDuration,
					CategoriesId: []string{testdata.FakeExistentCategoryId},
					GenresId:     []string{testdata.FakeExistentGenreId},
				},
			},
				nil,
			},
			wantErr: false,
		},
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockIdGen := mockId.NewMockIdGenerator(ctrl)
	mockRepo := mock.NewMockRepository(ctrl)
	mockLogger := mockLogger{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mockRepo.EXPECT().
					GetAll(tt.args.ctx).
					Return(
						tt.want.videos,
						nil,
					)
			}
			s := video.NewService(mockIdGen, mockRepo, mockLogger)
			got, err := s.List(tt.args.ctx)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.want.err.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.videos, got, "they should be equal")
			}
		})
	}
}

func Test_Service_Show(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := mock.NewMockRepository(ctrl)
	const (
		fakeDoesNotId = "fakeDoesNotId"
		fakeExistId   = "fakeExistId"
	)
	fakeErrorDoesNotId := fmt.Errorf("id %s: %w", fakeDoesNotId, logger.ErrNotFound)
	fakeErrorInternalApplication := fmt.Errorf("anyName: %w", logger.ErrInternalApplication)
	fakeVideo := video.Video{
		Id:           uuid.New().String(),
		Title:        fakeExistId,
		Description:  faker.Sentence(),
		YearLaunched: fakeYearLaunched,
		Opened:       false,
		Rating:       fakeRating,
		Duration:     fakeDuration,
	}
	type args struct {
		ctx context.Context
		id  string
	}
	type returns struct {
		video video.Video
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
			args: args{
				testdata.FakeCtx,
				"anyName",
			},
			want: returns{
				video.Video{},
				fakeErrorInternalApplication,
			},
			wantErr: true,
			setupMockR: func() {
				mockRepo.EXPECT().
					GetOne(testdata.FakeCtx, "anyName").
					Return(
						video.Video{},
						fakeErrorInternalApplication,
					)
			},
		},
		{
			name: "When id is not found",
			args: args{
				testdata.FakeCtx,
				fakeDoesNotId,
			},
			want: returns{
				video.Video{},
				fakeErrorDoesNotId,
			},
			wantErr: true,
			setupMockR: func() {
				mockRepo.EXPECT().
					GetOne(testdata.FakeCtx, fakeDoesNotId).
					Return(
						video.Video{},
						fakeErrorDoesNotId,
					)
			},
		},
		{
			name: "When id is found",
			args: args{
				testdata.FakeCtx,
				fakeExistId,
			},
			want: returns{
				fakeVideo,
				nil,
			},
			wantErr: false,
			setupMockR: func() {
				mockRepo.EXPECT().
					GetOne(testdata.FakeCtx, strings.TrimSpace(fakeExistId)).
					Return(
						fakeVideo,
						nil,
					)
			},
		},
	}
	mockIdGen := mockId.NewMockIdGenerator(ctrl)
	mockLogger := mockLogger{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMockR()
			s := video.NewService(mockIdGen, mockRepo, mockLogger)
			got, err := s.Show(tt.args.ctx, tt.args.id)
			if tt.wantErr {
				if assert.Error(t, err) {
					assert.EqualError(t, err, tt.want.err.Error(), "they should be equal")
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.video, got, "they should be equal")
			}
		})
	}
}

//func Test_Service_Update(t *testing.T) {
//	const (
//		fakeExistTitle        = "fakeExistTitle"
//		fakeDoesNotExistTitle = "fakeDoesNotExistTitle"
//	)
//	fakeTitle := faker.Name()
//	fakeBlankTitle := "     "
//	fakeDesc := faker.Sentence()
//	var fakeInvalidatedYearLaunched int16 = -999
//	false = false
//
//	type args struct {
//		ctx         context.Context
//		id          string
//		updateVideo video.UpdateVideo
//	}
//	tests := []struct {
//		name          string
//		args          args
//		expectedError error
//		wantErr       bool
//	}{
//		{
//			name: "When id is blank",
//			args: args{
//				testdata.FakeCtx,
//				"     ",
//				video.UpdateVideo{
//					Title: &fakeTitle,
//				},
//			},
//			expectedError: fmt.Errorf("'%s' field %w", "id", logger.ErrCouldNotBeEmpty),
//			wantErr:       true,
//		},
//		{
//			name: "When UpdateVideo is not provided",
//			args: args{
//				testdata.FakeCtx,
//				testdata.FakeId,
//				video.UpdateVideo{},
//			},
//			expectedError: fmt.Errorf("'%s' param %w", "updateVideo", logger.ErrCouldNotBeEmpty),
//			wantErr:       true,
//		},
//		{
//			name: "When the Title in UpdateVideo is blank",
//			args: args{
//				ctx: testdata.FakeCtx,
//				id:  testdata.FakeId,
//				updateVideo: video.UpdateVideo{
//					Title:        &fakeBlankTitle,
//					YearLaunched: &fakeYearLaunched,
//					Opened:       &false,
//					Rating:       &fakeRating,
//					Duration:     &fakeDuration,
//				}},
//			expectedError: fmt.Errorf("'Title' field %w", logger.ErrCouldNotBeEmpty),
//			wantErr:       true,
//		},
//		{
//			name: "When the YearLaunched in UpdateVideo is invalidated",
//			args: args{
//				ctx: testdata.FakeCtx,
//				id:  testdata.FakeId,
//				updateVideo: video.UpdateVideo{
//					Title:        &fakeTitle,
//					YearLaunched: &fakeInvalidatedYearLaunched,
//					Opened:       &false,
//					Rating:       &fakeRating,
//					Duration:     &fakeDuration,
//				}},
//			expectedError: fmt.Errorf("'YearLaunched' field %w", logger.ErrIsNotValidated),
//			wantErr:       true,
//		},
//		{
//			name: "When the YearLaunched in UpdateVideo is blank",
//			args: args{
//				ctx: testdata.FakeCtx,
//				id:  testdata.FakeId,
//				updateVideo: video.UpdateVideo{
//					Title:        &fakeTitle,
//					YearLaunched: nil,
//					Opened:       &false,
//					Rating:       &fakeRating,
//					Duration:     &fakeDuration,
//				}},
//			expectedError: nil,
//			wantErr:       false,
//		},
//		{
//			name: "When the Rating in UpdateVideo is blank",
//			args: args{
//				ctx: testdata.FakeCtx,
//				id:  testdata.FakeId,
//				updateVideo: video.UpdateVideo{
//					Title:        &fakeTitle,
//					YearLaunched: &fakeYearLaunched,
//					Opened:       &false,
//					Rating:       nil,
//					Duration:     &fakeDuration,
//				}},
//			expectedError: nil,
//			wantErr:       false,
//		},
//		{
//			name: "When the Duration in UpdateVideo is blank",
//			args: args{
//				ctx: testdata.FakeCtx,
//				id:  testdata.FakeId,
//				updateVideo: video.UpdateVideo{
//					Title:        &fakeTitle,
//					YearLaunched: &fakeYearLaunched,
//					Opened:       &false,
//					Rating:       &fakeRating,
//					Duration:     nil,
//				},
//			},
//			expectedError: nil,
//			wantErr:       false,
//		},
//		{
//			name: "When UpdateVideo is with wrong CategoriesId",
//			args: args{
//				testdata.FakeCtx,
//				testdata.FakeId,
//				video.UpdateVideo{
//					Title:        &fakeTitle,
//					Description:  &fakeDesc,
//					YearLaunched: &fakeYearLaunched,
//					Opened:       &false,
//					Rating:       &fakeRating,
//					Duration:     &fakeDuration,
//					GenresId:     []string{testdata.FakeExistentGenreId},
//					CategoriesId: []string{testdata.FakeNonExistentCategoryId},
//				},
//			},
//			expectedError: logger.ErrIsRequired,
//			wantErr:       true,
//		},
//		{
//			name: "When UpdateVideo is without CategoriesId and GenresId",
//			args: args{
//				testdata.FakeCtx,
//				testdata.FakeId,
//				video.UpdateVideo{
//					Title:        &fakeTitle,
//					Description:  &fakeDesc,
//					YearLaunched: &fakeYearLaunched,
//					Opened:       &false,
//					Rating:       &fakeRating,
//					Duration:     &fakeDuration,
//				},
//			},
//			expectedError: fmt.Errorf("'Categories' field %w", logger.ErrIsRequired),
//			wantErr:       true,
//		},
//		{
//			name: "When id is not found",
//			args: args{
//				testdata.FakeCtx,
//				testdata.FakeId,
//				video.UpdateVideo{
//					Title:        &fakeTitle,
//					YearLaunched: &fakeYearLaunched,
//					Opened:       &false,
//					Rating:       &fakeRating,
//					Duration:     &fakeDuration,
//					GenresId:     []string{testdata.FakeExistentGenreId},
//					CategoriesId: []string{testdata.FakeExistentCategoryId},
//				},
//			},
//			expectedError: fmt.Errorf("%s: %w", fakeDoesNotExistTitle, logger.ErrNotFound),
//			wantErr:       true,
//		},
//		{
//			name: "When id is found and UpdateVideo is right",
//			args: args{
//				testdata.FakeCtx,
//				fakeExistTitle,
//				video.UpdateVideo{
//					Title:        &fakeTitle,
//					YearLaunched: &fakeYearLaunched,
//					Opened:       &false,
//					Rating:       &fakeRating,
//					Duration:     &fakeDuration,
//					GenresId:     []string{testdata.FakeExistentGenreId},
//					CategoriesId: []string{testdata.FakeExistentCategoryId},
//				},
//			},
//			expectedError: nil,
//			wantErr:       false,
//		},
//	}
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	mockIdGen := mockId.NewMockIdGenerator(ctrl)
//	mockRepo := mock.NewMockRepository(ctrl)
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			switch tt.name {
//			case "When the Title in Category already exists",
//				"When the YearLaunched in UpdateVideo is blank",
//				"When UpdateVideo is without CategoriesId and GenresId",
//				"When id is not found",
//				"When id is found and UpdateVideo is right":
//				mockRepo.EXPECT().
//					UpdateOne(tt.args.ctx, tt.args.id, tt.args.updateVideo).
//					Return(tt.expectedError)
//			}
//			s := video.NewService(mockIdGen, mockRepo)
//			fmt.Println(tt.name)
//			if tt.name == "When the Rating in UpdateVideo is blank" && tt.args.updateVideo.Title != nil {
//				fmt.Println("Title-1: ", *tt.args.updateVideo.Title)
//			}
//			fmt.Printf("1: %v %v %#v\n", tt.args.ctx, tt.args.id, tt.args.updateVideo)
//			err := s.Update(tt.args.ctx, tt.args.id, tt.args.updateVideo)
//			if tt.wantErr {
//				if assert.Error(t, err) {
//					assert.EqualError(t, err, tt.expectedError.Error(), "they should be equal")
//				}
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
