package errors

import (
	"errors"
)

var (
	ErrNotFound        = errors.New("entity not found")
	ErrDuplicateEntity = errors.New("entity already exists")
	ErrInvalidStatus   = errors.New("invalid status")
)
