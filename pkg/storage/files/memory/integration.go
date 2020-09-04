package memory

import (
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/afero"

	"github.com/selmison/code-micro-videos/pkg/storage/files"
)

type FileSeed struct {
	Repo                     files.Repository
	FakeAfero                *afero.Afero
	FakeVideoIDExists        uuid.UUID
	FakeVideoIDDoesNotExist  uuid.UUID
	FakeVideoFileExists      afero.File
	FakeVideoFileNameExists  string
	FakeFileNameDoesNotExist string
	FakeTmpFile              afero.File
	FakeTmpDir               string
	FakeTmpFileName          string
	FakeTmpFilePath          string
}

func SetupFileTestCase() (*FileSeed, func(t *testing.T), error) {
	repo := NewRepository()
	fakeAfs := repo.Afs
	fakeVideoIDExists := uuid.New()
	fakeVideoIDDoesNotExist := uuid.New()
	if err := fakeAfs.Mkdir(fakeVideoIDExists.String(), 0755); err != nil {
		return nil, nil, fmt.Errorf("test: failed to make video directory: %v\n", err)
	}
	fakeFileExists, err := fakeAfs.TempFile(fakeVideoIDExists.String(), "")
	defer func() {
		if err := fakeFileExists.Close(); err != nil {
			log.Printf("test: could not close file: %v", err)
		}
	}()
	fakeFileStatExists, err := fakeFileExists.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("test: failed to get file info: %v\n", err)
	}
	fakeFileNameExists := fakeFileStatExists.Name()
	fakeFileNameDoesNotExist := "FakeFileNameDoesNotExist"
	fakeData := make([]byte, 20)
	if _, err := rand.Read(fakeData); err != nil {
		return nil, nil, fmt.Errorf("test: failed to generate a random Data: %v\n", err)
	}
	if err := fakeAfs.WriteFile(fakeFileExists.Name(), fakeData, 0644); err != nil {
		return nil, nil, fmt.Errorf("test: failed to write a new file: %v\n", err)
	}
	fakeTmpDir := "tmpDir"
	if err := fakeAfs.Mkdir(fakeTmpDir, 0755); err != nil {
		return nil, nil, fmt.Errorf("test: failed to make temp directory: %v\n", err)
	}
	fakeTmpFile, err := fakeAfs.TempFile(fakeTmpDir, "")
	defer func() {
		if err := fakeTmpFile.Close(); err != nil {
			log.Printf("test: could not close file: %v", err)
		}
	}()
	fakeTmpFileStatExists, err := fakeFileExists.Stat()
	if err != nil {
		return nil, nil, fmt.Errorf("test: failed to get file info: %v\n", err)
	}
	fakeTmpFileName := fakeTmpFileStatExists.Name()
	fakeTmpFilePath := fakeTmpFile.Name()
	if err := fakeAfs.WriteFile(fakeTmpFile.Name(), fakeData, 0644); err != nil {
		return nil, nil, fmt.Errorf("test: failed to write a new file: %v\n", err)
	}
	s := &FileSeed{
		Repo:                     repo,
		FakeAfero:                fakeAfs,
		FakeVideoIDExists:        fakeVideoIDExists,
		FakeVideoIDDoesNotExist:  fakeVideoIDDoesNotExist,
		FakeVideoFileExists:      fakeFileExists,
		FakeVideoFileNameExists:  fakeFileNameExists,
		FakeFileNameDoesNotExist: fakeFileNameDoesNotExist,
		FakeTmpDir:               fakeTmpDir,
		FakeTmpFile:              fakeTmpFile,
		FakeTmpFileName:          fakeTmpFileName,
		FakeTmpFilePath:          fakeTmpFilePath,
	}
	return s, func(t *testing.T) {
		defer func() {
			if err := fakeFileExists.Close(); err != nil {
				t.Errorf("test: failed to close file: %v", err)
			}
		}()
	}, nil
}
