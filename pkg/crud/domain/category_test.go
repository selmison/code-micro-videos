package domain

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

func TestCategory_Validate(t *testing.T) {
	type fields struct {
		Id          string
		Name        string
		Description string
		Genres      []Genre
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
				Name: strings.ToLower(fakeName),
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
			name: "when name is not lower",
			fields: fields{
				Id:   fakeUUID,
				Name: fakeName,
			},
			expectedError: fmt.Sprintf("lowercase 'Name' field %v", logger.ErrIsRequired),
			wantErr:       true,
		},
		{
			name: "when everything is right",
			fields: fields{
				Id:   fakeUUID,
				Name: strings.ToLower(fakeName),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Category{
				Id:          tt.fields.Id,
				Name:        tt.fields.Name,
				Description: tt.fields.Description,
				Genres:      tt.fields.Genres,
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
