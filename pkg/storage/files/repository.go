package files

import (
	"github.com/google/uuid"
)

type Repository interface {
	Exists(videoID uuid.UUID, fileName string) (bool, error)
	GetFileFromVideo(videoUUID uuid.UUID, filename string) ([]byte, error)
	SaveFileToVideo(videoID uuid.UUID, fileName string, fileData []byte) error
	UpdateFileToVideo(videoID uuid.UUID, fileName string, fileData []byte) (bool, error)
}
