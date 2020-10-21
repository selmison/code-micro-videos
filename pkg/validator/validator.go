package validator

import (
	"fmt"
)

type Validator interface {
	Validate() error
	IsValidated() bool
}

type check struct {
	hasError bool
	err      error
}

func NewCheck(hasError bool, err error) check {
	return check{hasError, err}
}

func CheckAll(checks ...check) error {
	for _, check := range checks {
		if check.hasError {
			return check.err
		}
	}
	return nil
}

type validator struct {
	err error
}

func NewValidator() Validator {
	return &validator{}
}

func (v *validator) MustBeGreaterThan(high, value int) bool {
	if v.err != nil {
		return false
	}
	if value <= high {
		v.err = fmt.Errorf("must be Greater than %d", high)
		return false
	}
	return true
}

func (v *validator) MustBeNotEmpty(value string) bool {
	if v.err != nil {
		return false
	}
	if value == "" {
		v.err = fmt.Errorf("must not be empty")
		return false
	}
	return true
}

func (v *validator) IsValidated() bool {
	return v.err != nil
}

func (v *validator) Validate() error {
	return v.err
}
