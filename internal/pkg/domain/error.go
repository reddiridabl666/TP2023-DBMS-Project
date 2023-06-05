package domain

import "errors"

//easyjson:json
type ErrorMessage struct {
	Message string
}

var (
	ErrUniqueViolation = errors.New("such object already exists")
	ErrNotFound        = errors.New("no such object")
)
