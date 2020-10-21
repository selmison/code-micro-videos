package logger

import (
	"errors"
	"fmt"
)

type Errorer interface {
	Error() error
}

type ResultError struct {
	ErrMsg string
	Err    error
}

var (
	ErrInvalidedLimit      = errors.New("limit should be positive number")
	ErrNotFound            = errors.New("not found")
	ErrInternalApplication = errors.New("internal application error")
	ErrIsNotValidated      = errors.New("is not validated")
	ErrIsRequired          = errors.New("is required")
	ErrCouldNotBeEmpty     = errors.New("could not be empty")
	ErrAlreadyExists       = errors.New("already exists")
)

func (e *ResultError) Unwrap() error { return e.Err }

func (e *ResultError) Error() string {
	return fmt.Sprint(e.ErrMsg, " ", e.Err.Error())
}
