package domain

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/selmison/code-micro-videos/pkg/logger"
)

var (
	fakeUUID = uuid.New().String()
	fakeName = faker.FirstName()
)

func TestGenreValidatable_Validate(t *testing.T) {

	type fields struct {
		Id         string
		Name       string
		Categories []CategoryValidatable
	}
	tests := []struct {
		name    string
		fields  fields
		returns string
		wantErr bool
	}{
		{
			name: "when id is blank",
			fields: fields{
				Id:   "     ",
				Name: strings.ToLower(fakeName),
			},
			returns: fmt.Sprintf("'Id' field %v", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name: "when name is blank",
			fields: fields{
				Id:   fakeUUID,
				Name: "   ",
			},
			returns: fmt.Sprintf("'Name' field %v", logger.ErrIsRequired),
			wantErr: true,
		},
		{
			name: "when name is not lower",
			fields: fields{
				Id:   fakeUUID,
				Name: fakeName,
			},
			returns: fmt.Sprintf("lowercase 'Name' field %v", logger.ErrIsRequired),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GenreValidatable{
				Id:         tt.fields.Id,
				Name:       tt.fields.Name,
				Categories: tt.fields.Categories,
			}
			err := g.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error: %v, wantErr %v", err, tt.wantErr)
			}
			assert.EqualError(t, err, tt.returns, "they should be equal")
		})
	}
}

func TestGenreValidatable_mapToGenre(t *testing.T) {
	type fields struct {
		Id         string
		Name       string
		Categories []CategoryValidatable
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Genre
		returns string
		wantErr bool
	}{
		{
			name: "when name is not lower",
			fields: fields{
				Id:   fakeUUID,
				Name: fakeName,
			},
			returns: fmt.Sprintf("lowercase 'Name' field %v", logger.ErrIsRequired),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GenreValidatable{
				Id:         tt.fields.Id,
				Name:       tt.fields.Name,
				Categories: tt.fields.Categories,
			}
			got, err := g.mapToGenre()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error: %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapToGenre() = %v, want %v", got, tt.want)
			}
			assert.EqualError(t, err, tt.returns, "they should be equal")
		})
	}
}

func TestGenre_MapToGenreValidatable(t *testing.T) {
	type fields struct {
		id         string
		name       string
		categories []Category
	}
	tests := []struct {
		name   string
		fields fields
		want   *GenreValidatable
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Genre{
				id:         tt.fields.id,
				name:       tt.fields.name,
				categories: tt.fields.categories,
			}
			if got := g.MapToGenreValidatable(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapToGenreValidatable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenre_Name(t *testing.T) {
	type fields struct {
		id         string
		name       string
		categories []Category
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Genre{
				id:         tt.fields.id,
				name:       tt.fields.name,
				categories: tt.fields.categories,
			}
			if got := g.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGenre(t *testing.T) {
	type args struct {
		fields GenreValidatable
	}
	tests := []struct {
		name    string
		args    args
		want    *Genre
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGenre(tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGenre() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGenre() got = %v, want %v", got, tt.want)
			}
		})
	}
}
