package domain

import "errors"

var (
	ErrBadRequest       = errors.New("bad request")
	ErrToolFailed       = errors.New("tool failed")
	ErrModelUnavailable = errors.New("model unavailable")
	ErrNoContextFound   = errors.New("no context found")
)
