package logger

import (
	"errors"
)

var (
	ErrInvalidedLimit      = errors.New("limit should be positive number")
	ErrNotFound            = errors.New("not found")
	ErrInternalApplication = errors.New("internal application error")
	ErrIsNotValidated      = errors.New("is not validated")
	ErrIsRequired          = errors.New("is required")
	ErrAlreadyExists       = errors.New("already exists")
)
