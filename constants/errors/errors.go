package errors

import "errors"

var (
	ErrCommandNotFound    = errors.New("command not found")
	ErrEventNotFound      = errors.New("event not found")
	ErrInvalidPayloadType = errors.New("received an invalid payload type")
)
