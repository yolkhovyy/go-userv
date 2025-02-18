package usergraphql

import "errors"

var (
	ErrDeleteFailure        = errors.New("delete failure")
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)
