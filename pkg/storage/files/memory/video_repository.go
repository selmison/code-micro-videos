//+build integration

package memory

import (
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/afero"
)

type repository struct {
	Afs *afero.Afero
}

func NewRepository() *repository {
	fs := afero.NewMemMapFs()
	r := &repository{
		Afs: &afero.Afero{Fs: fs},
	}
	return r
}

func (r *repository) Exists(videoID uuid.UUID, fileName string) (bool, error) {
	var filePath string
	if videoID == (uuid.UUID{}) {
		filePath = fileName
	} else {
		filePath = fmt.Sprintf("%s%c%s", videoID, os.PathSeparator, fileName)
	}
	exists, err := r.Afs.Exists(filePath)
	if err != nil {
		return false, fmt.Errorf("could not verify if file exists: %v", err)
	}
	return exists, nil
}

func (r *repository) GetFileFromVideo(videoID uuid.UUID, fileName string) ([]byte, error) {
	filePath := fmt.Sprintf("%s%c%s", videoID, os.PathSeparator, fileName)
	return r.Afs.ReadFile(filePath)
}

func (r *repository) SaveFileToVideo(videoID uuid.UUID, fileName string, fileData io.Reader) error {
	filePath := fmt.Sprintf("%s%c%s", videoID, os.PathSeparator, fileName)
	if fileData != nil {
		if err := r.Afs.WriteReader(filePath, fileData); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) UpdateFileToVideo(videoID uuid.UUID, fileName string, fileData io.Reader) (bool, error) {
	exists, err := r.Exists(videoID, fileName)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}
	if err := r.SaveFileToVideo(videoID, fileName, fileData); err != nil {
		return false, err
	}
	return true, err
}
