package domain

import "errors"

var (
	ErrAlreadyExists = errors.New("such object already exists")
	ErrNotFound      = errors.New("no such object")
)
