package domain

import (
	"fmt"
)

type Validator interface {
	Validate() error
	IsValidated() bool
}

type validator struct {
	err error
}

func (e *Event) Validate() error {
	return Check(
		Cf(e.Id <= 0, "Expected ID <= 0, got %d.", e.Id),
		Cf(e.Start.UnixNano() > e.End.UnixNano(), "Expected start < end, got %s >= %s.", e.Start, e.End),
	)
}

type C struct {
	check bool
	err   error
}

func Cf(chk bool, errMsg string, params ...interface{}) C {
	return C{chk, fmt.Errorf(errMsg, params...)}
}

func Check(args ...C) error {
	for _, c := range args {
		if !c.check {
			return c.err
		}
	}
	return nil
}

//func (v *validator) MustBeGreaterThan(high, value int) bool {
//	if v.err != nil {
//		return false
//	}
//	if value <= high {
//		v.err = fmt.Errorf("must be Greater than %d", high)
//		return false
//	}
//	return true
//}
//
//func (v *validator) MustBeBefore(high, value time.Time) bool {
//	if v.err != nil {
//		return false
//	}
//	if value.After(high) {
//		v.err = fmt.Errorf("must be Before than %v", high)
//		return false
//	}
//	return true
//}
//
//func (v *validator) MustBeNotEmpty(value string) bool {
//	if v.err != nil {
//		return false
//	}
//	if value == "" {
//		v.err = fmt.Errorf("must not be empty")
//		return false
//	}
//	return true
//}

func (v *validator) IsValid() bool {
	return v.err != nil
}

func (v *validator) Error() string {
	return v.err.Error()
}
