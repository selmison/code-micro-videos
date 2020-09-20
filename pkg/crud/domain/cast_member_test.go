package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

const fakeNotValidatedCastMemberType = 111

func TestCastMemberType_String(t *testing.T) {
	tests := []struct {
		name string
		c    CastMemberType
		want string
	}{
		{
			name: "when CastMemberType is Director",
			c:    Director,
			want: "Director",
		},
		{
			name: "when CastMemberType is Actor",
			c:    Actor,
			want: "Actor",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.c.String(), tt.want, "they should be equal")
		})
	}
}

func TestCastMemberType_Validate(t *testing.T) {
	tests := []struct {
		name          string
		c             CastMemberType
		expectedError string
		wantErr       bool
	}{
		{
			name:          "when CastMemberType is not validated",
			c:             111,
			expectedError: fmt.Sprintf("cast member type %v", logger.ErrIsNotValidated),
			wantErr:       true,
		},
		{
			name:    "when CastMemberType is Director",
			c:       Director,
			wantErr: false,
		},
		{
			name:    "when CastMemberType is Actor",
			c:       Actor,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.Validate()
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

func TestCastMember_Validate(t *testing.T) {
	type fields struct {
		Id   string
		Name string
		Type CastMemberType
	}
	tests := []struct {
		name          string
		fields        fields
		expectedError string
		wantErr       bool
	}{
		{
			name: "when id is blank",
			fields: fields{
				Id:   "     ",
				Name: fakeName,
			},
			expectedError: fmt.Sprintf("'Id' field %v", logger.ErrIsRequired),
			wantErr:       true,
		},
		{
			name: "when name is blank",
			fields: fields{
				Id:   fakeUUID,
				Name: "   ",
			},
			expectedError: fmt.Sprintf("'Name' field %v", logger.ErrIsRequired),
			wantErr:       true,
		},
		{
			name: "when rating is not validated",
			fields: fields{
				Id:   fakeUUID,
				Name: fakeName,
				Type: fakeNotValidatedCastMemberType,
			},
			expectedError: fmt.Sprintf("cast member type %v", logger.ErrIsNotValidated),
			wantErr:       true,
		},
		{
			name: "when everything is right",
			fields: fields{
				Id:   fakeUUID,
				Name: fakeName,
				Type: Director,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CastMember{
				Id:   tt.fields.Id,
				Name: tt.fields.Name,
				Type: tt.fields.Type,
			}
			err := c.Validate()
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
