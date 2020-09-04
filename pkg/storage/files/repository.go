package files

import (
	"io"

	"github.com/google/uuid"
)

type Repository interface {
	Exists(videoID uuid.UUID, fileName string) (bool, error)
	GetFileFromVideo(videoID uuid.UUID, fileName string) ([]byte, error)
	SaveFileToVideo(videoID uuid.UUID, fileName string, fileData io.Reader) error
	UpdateFileToVideo(videoID uuid.UUID, fileName string, fileData io.Reader) (bool, error)
}
