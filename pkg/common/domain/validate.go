package domain

type Validator interface {
	Validate() error
	IsValidated() bool
}
