package repository

import "errors"

// ErrNotFound is returned when a single-row query finds no matching row.
var ErrNotFound = errors.New("repository: not found")
