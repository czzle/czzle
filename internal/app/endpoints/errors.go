package endpoints

import "errors"

var (
	ErrTypeAssertion     = errors.New("endpoint type assertion")
	ErrServiceNotDefined = errors.New("endpoint service not defined")
)
