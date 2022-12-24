package httpErrors

import "errors"

var (
	ErrBadResult  = errors.New("bad result")
	ErrBadRequest = errors.New("bad request")
	ErrWrite      = errors.New("write error")
	ErrGetValue   = errors.New("cant get value for this range")
)
