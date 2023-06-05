package domain

import "errors"

var (
	ErrUniqueViolation = errors.New("such object already exists")
	ErrNotFound        = errors.New("no such object")
)
