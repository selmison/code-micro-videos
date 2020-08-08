package logger

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidedLimit      = fmt.Errorf("limit should be positive number")
	ErrNotFound            = errors.New("not found")
	ErrInternalApplication = errors.New("internal application error")
	ErrIsNotValidated      = fmt.Errorf("is not validated")
	ErrIsRequired          = fmt.Errorf("is required")
	ErrAlreadyExists       = fmt.Errorf("already exists")
)
