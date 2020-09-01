package files

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/afero"
)

func Test_repository_Exists(t *testing.T) {
	seed, teardownTestCase, err := setupTestCase()
	if err != nil {
		t.Fatalf("test: failed to setup test case: %v\n", err)
	}
	defer teardownTestCase(t)
	type args struct {
		videoID  uuid.UUID
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
		err     error
	}{
		{
			name: "when video does not exist",
			args: args{
				videoID:  seed.fakeVideoIDDoesNotExist,
				fileName: seed.fakeFileNameExists,
			},
			want:    false,
			wantErr: false,
			err:     nil,
		},
		{
			name: "when file does not exist",
			args: args{
				videoID:  seed.fakeVideoIDExists,
				fileName: seed.fakeFileNameDoesNotExist,
			},
			want:    false,
			wantErr: false,
			err:     nil,
		},
		{
			name: "whe file exists",
			args: args{
				videoID:  seed.fakeVideoIDExists,
				fileName: seed.fakeFileNameExists,
			},
			want:    true,
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := seed.repo.Exists(tt.args.videoID, tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Exists() error = %v, wantErr %v\n", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("Exists() got: %v, want: %v\n", got, tt.want)
			}
		})
	}
}

func Test_repository_GetFileFromVideo(t *testing.T) {
	seed, teardownTestCase, err := setupTestCase()
	if err != nil {
		t.Fatalf("test: failed to setup test case: %v\n", err)
	}
	sb, err := afero.ReadFile(seed.repo.Afs.Fs, seed.fakeFileExists.Name())
	if err != nil {
		t.Fatalf("test: could not read file: %v", err)
	}
	defer teardownTestCase(t)
	type args struct {
		videoID  uuid.UUID
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "when video does not exist",
			args: args{
				videoID:  seed.fakeVideoIDDoesNotExist,
				fileName: seed.fakeFileNameExists,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when file does not exist",
			args: args{
				videoID:  seed.fakeVideoIDExists,
				fileName: seed.fakeFileNameDoesNotExist,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "whe file exists",
			args: args{
				videoID:  seed.fakeVideoIDExists,
				fileName: seed.fakeFileNameExists,
			},
			want:    sb,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := seed.repo.GetFileFromVideo(tt.args.videoID, tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileFromVideo() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFileFromVideo() got: %v, want: %v", got, tt.want)
			}
		})
	}
}

func Test_repository_SaveFileToVideo(t *testing.T) {
	seed, teardownTestCase, err := setupTestCase()
	if err != nil {
		t.Fatalf("test: failed to setup test case: %v\n", err)
	}
	defer teardownTestCase(t)
	type args struct {
		videoID  uuid.UUID
		fileName string
		fileData []byte
	}
	tests := []struct {
		name    string
		args    args
		err     error
		wantErr bool
	}{
		{
			name: "when video does not exist",
			args: args{
				videoID:  seed.fakeVideoIDDoesNotExist,
				fileName: seed.fakeFileNameExists,
			},
			err:     nil,
			wantErr: false,
		},
		{
			name: "when file does not exist",
			args: args{
				videoID:  seed.fakeVideoIDExists,
				fileName: seed.fakeFileNameDoesNotExist,
			},
			err:     nil,
			wantErr: false,
		},
		{
			name: "whe file exists",
			args: args{
				videoID:  seed.fakeVideoIDExists,
				fileName: seed.fakeFileNameExists,
			},
			err:     nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := seed.repo.SaveFileToVideo(tt.args.videoID, tt.args.fileName, tt.args.fileData); (err != nil) != tt.wantErr {
				t.Fatalf("SaveFileToVideo() error = %v\n, wantErr %v", err, tt.wantErr)
			}
			filePath := fmt.Sprintf("%s%c%s", tt.args.videoID, os.PathSeparator, tt.args.fileName)
			exists, err := afero.Exists(seed.repo.Afs.Fs, filePath)
			if err != nil {
				t.Fatalf("test: could not verify if file exist: %v\n", err)
			}
			if !exists {
				t.Fatalf("test: file '%s' doesn't exist\n", filePath)
			}
		})
	}
}

func Test_repository_UpdateFileToVideo(t *testing.T) {
	seed, teardownTestCase, err := setupTestCase()
	if err != nil {
		t.Fatalf("test: failed to setup test case: %v\n", err)
	}
	defer teardownTestCase(t)
	type fields struct {
		Afs     *afero.Afero
		tempDir string
	}
	type args struct {
		videoID  uuid.UUID
		fileName string
		fileData []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "when file does not exist",
			args: args{
				videoID:  seed.fakeVideoIDExists,
				fileName: seed.fakeFileNameDoesNotExist,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "whe file exists",
			args: args{
				videoID:  seed.fakeVideoIDExists,
				fileName: seed.fakeFileNameExists,
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := seed.repo.UpdateFileToVideo(tt.args.videoID, tt.args.fileName, tt.args.fileData)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateFileToVideo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateFileToVideo() got = %v, want %v", got, tt.want)
			}
		})
	}
}
