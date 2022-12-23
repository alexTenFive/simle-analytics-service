package httpErrors

import "errors"

var (
	ErrBadResult  error = errors.New("bad result")
	ErrBadRequest error = errors.New("bad request")
	ErrWrite      error = errors.New("write error")
	ErrGetValue   error = errors.New("cant get value for this range")
)
