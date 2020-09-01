package files

import (
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/afero"
)

type seed struct {
	repo                     *repository
	fakeVideoIDExists        uuid.UUID
	fakeVideoIDDoesNotExist  uuid.UUID
	fakeFileExists           afero.File
	fakeFileNameExists       string
	fakeFileNameDoesNotExist string
}

func setupTestCase() (*seed, func(t *testing.T), error) {
	fakeAfs := NewRepository()
	fakeVideoIDExists := uuid.New()
	fakeVideoIDDoesNotExist := uuid.New()
	if err := fakeAfs.Afs.Mkdir(fakeVideoIDExists.String(), 0755); err != nil {
		return nil, nil, fmt.Errorf("test: failed to make video directory: %v\n", err)
	}
	fakeFileExists, err := fakeAfs.Afs.TempFile(fakeVideoIDExists.String(), "")
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
	fakeFileNameDoesNotExist := "fakeFileNameDoesNotExist"
	fakeData := make([]byte, 2<<20)
	if _, err := rand.Read(fakeData); err != nil {
		return nil, nil, fmt.Errorf("test: failed to generate a random Data: %v\n", err)
	}
	if err := fakeAfs.Afs.WriteFile(fakeFileExists.Name(), fakeData, 0644); err != nil {
		return nil, nil, fmt.Errorf("test: failed to write a new file: %v\n", err)
	}
	s := &seed{
		repo:                     fakeAfs,
		fakeVideoIDExists:        fakeVideoIDExists,
		fakeVideoIDDoesNotExist:  fakeVideoIDDoesNotExist,
		fakeFileExists:           fakeFileExists,
		fakeFileNameExists:       fakeFileNameExists,
		fakeFileNameDoesNotExist: fakeFileNameDoesNotExist,
	}
	return s, func(t *testing.T) {
		defer func() {
			if err := fakeFileExists.Close(); err != nil {
				t.Errorf("test: failed to close file: %v", err)
			}
		}()
	}, nil
}
