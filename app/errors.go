package app

import "errors"

var (
	ErrBadRequest = errors.New("bad request")
	ErrNotFound   = errors.New("not found")
	ErrCollision  = errors.New("failed to create short link due to collision")
)
