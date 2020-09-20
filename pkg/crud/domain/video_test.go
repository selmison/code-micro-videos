package domain

import (
	"fmt"
	"mime/multipart"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

const (
	fakeOpened = false
)

var (
	fakeTitle                    = faker.Name()
	fakeDesc                     = faker.Sentence()
	fakeYearLaunched       int16 = 2020
	fakeDuration           int16 = 90
	fakeRating                   = TwelveRating
	fakeNotValidatedRating       = VideoRating(111)
)

func TestVideoRating_String(t *testing.T) {
	tests := []struct {
		name string
		v    VideoRating
		want string
	}{
		{
			name: "When VideoRating is FreeRating",
			v:    FreeRating,
			want: "Free",
		},
		{
			name: "When VideoRating is TenRating",
			v:    TenRating,
			want: "10",
		},
		{
			name: "When VideoRating is TwelveRating",
			v:    TwelveRating,
			want: "12",
		},
		{
			name: "When VideoRating is FourteenRating",
			v:    FourteenRating,
			want: "14",
		},
		{
			name: "When VideoRating is SixteenRating",
			v:    SixteenRating,
			want: "16",
		},
		{
			name: "When VideoRating is EighteenRating",
			v:    EighteenRating,
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
		v             VideoRating
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
			v:       FreeRating,
			wantErr: false,
		},
		{
			name:    "When VideoRating is TenRating",
			v:       TenRating,
			wantErr: false,
		},
		{
			name:    "When VideoRating is TwelveRating",
			v:       TwelveRating,
			wantErr: false,
		},
		{
			name:    "When VideoRating is FourteenRating",
			v:       FourteenRating,
			wantErr: false,
		},
		{
			name:    "When VideoRating is SixteenRating",
			v:       SixteenRating,
			wantErr: false,
		},
		{
			name:    "When VideoRating is EighteenRating",
			v:       EighteenRating,
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

func TestVideo_Validate(t *testing.T) {
	type fields struct {
		Id               string
		Title            string
		Description      string
		YearLaunched     *int16
		Opened           bool
		Rating           *VideoRating
		Duration         *int16
		Categories       []Category
		Genres           []Genre
		VideoFileHandler *multipart.FileHeader
	}
	tests := []struct {
		name          string
		fields        fields
		expectedError string
		wantErr       bool
	}{
		{
			name: "When the Id is blank",
			fields: fields{
				Id:           "    ",
				Title:        fakeTitle,
				YearLaunched: &fakeYearLaunched,
				Opened:       fakeOpened,
				Rating:       &fakeRating,
				Duration:     &fakeDuration,
			},
			expectedError: fmt.Sprintf("'Id' field %v", logger.ErrIsRequired),
			wantErr:       true,
		},
		{
			name: "When the Title is blank",
			fields: fields{
				Id:           fakeUUID,
				Title:        "    ",
				YearLaunched: &fakeYearLaunched,
				Opened:       fakeOpened,
				Rating:       &fakeRating,
				Duration:     &fakeDuration,
			},
			expectedError: fmt.Sprintf("'Title' field %v", logger.ErrIsRequired),
			wantErr:       true,
		},
		{
			name: "When the YearLaunched is blank",
			fields: fields{
				Id:           fakeUUID,
				Title:        fakeTitle,
				YearLaunched: nil,
				Opened:       fakeOpened,
				Rating:       &fakeRating,
				Duration:     &fakeDuration,
			},
			expectedError: fmt.Sprintf("'YearLaunched' field %v", logger.ErrIsRequired),
			wantErr:       true,
		},
		{
			name: "When the Rating is blank",
			fields: fields{
				Id:           fakeUUID,
				Title:        fakeTitle,
				YearLaunched: &fakeYearLaunched,
				Opened:       fakeOpened,
				Rating:       nil,
				Duration:     &fakeDuration,
			},
			expectedError: fmt.Sprintf("'Rating' field %v", logger.ErrIsRequired),
			wantErr:       true,
		},
		{
			name: "When the Duration is blank",
			fields: fields{
				Id:           fakeUUID,
				Title:        fakeTitle,
				YearLaunched: &fakeYearLaunched,
				Opened:       fakeOpened,
				Rating:       &fakeRating,
				Duration:     nil,
			},
			expectedError: fmt.Sprintf("'Duration' field %v", logger.ErrIsRequired),
			wantErr:       true,
		},
		{
			name: "When Video is without genres and categories",
			fields: fields{
				Id:           fakeUUID,
				Title:        fakeTitle,
				Description:  fakeDesc,
				YearLaunched: &fakeYearLaunched,
				Opened:       fakeOpened,
				Rating:       &fakeRating,
				Duration:     &fakeDuration,
			},
			expectedError: fmt.Sprintf("'Categories' field %v", logger.ErrIsRequired),
			wantErr:       true,
		},
		{
			name: "When VideoRating is not validated",
			fields: fields{
				Id:           faker.UUIDHyphenated(),
				Title:        fakeTitle,
				Description:  fakeDesc,
				YearLaunched: &fakeYearLaunched,
				Opened:       fakeOpened,
				Rating:       &fakeNotValidatedRating,
				Duration:     &fakeDuration,
				Genres: []Genre{
					{
						Id:   faker.UUIDHyphenated(),
						Name: faker.FirstName(),
					},
				},
				Categories: []Category{
					{
						Id:   faker.UUIDHyphenated(),
						Name: faker.FirstName(),
					},
				},
			},
			expectedError: fmt.Sprintf("video rating %v", logger.ErrIsNotValidated),
			wantErr:       true,
		},
		{
			name: "When everything is right",
			fields: fields{
				Id:           faker.UUIDHyphenated(),
				Title:        fakeTitle,
				Description:  fakeDesc,
				YearLaunched: &fakeYearLaunched,
				Opened:       fakeOpened,
				Rating:       &fakeRating,
				Duration:     &fakeDuration,
				Genres: []Genre{
					{
						Id:   faker.UUIDHyphenated(),
						Name: faker.FirstName(),
					},
				},
				Categories: []Category{
					{
						Id:   faker.UUIDHyphenated(),
						Name: faker.FirstName(),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Video{
				Id:               tt.fields.Id,
				Title:            tt.fields.Title,
				Description:      tt.fields.Description,
				YearLaunched:     tt.fields.YearLaunched,
				Opened:           tt.fields.Opened,
				Rating:           tt.fields.Rating,
				Duration:         tt.fields.Duration,
				Categories:       tt.fields.Categories,
				Genres:           tt.fields.Genres,
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
