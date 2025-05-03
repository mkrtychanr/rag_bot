package screen

import "errors"

var (
	ErrUnknownScreen = errors.New("unknown screen")
	ErrEmptyPayload  = errors.New("empty payload")
	ErrWrongType     = errors.New("wrong type")
)
