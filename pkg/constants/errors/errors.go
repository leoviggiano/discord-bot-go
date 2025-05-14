package errors

import "errors"

var (
	ErrCommandExpired     = errors.New("command expired")
	ErrCommandNotFound    = errors.New("command not found")
	ErrEventNotFound      = errors.New("event not found")
	ErrInvalidPayloadType = errors.New("received an invalid payload type")
	ErrUserNotFound       = errors.New("user not found")

	ErrInvalidAmount = errors.New("invalid amount")

	ErrUserAlreadyInBattle  = errors.New("user already in battle")
	ErrUserNotFoundInBattle = errors.New("user not found in battle data")
	ErrSkillNotFound        = errors.New("skill not found")
	ErrBuffEmptyAttribute   = errors.New("buff attribute is empty")
	ErrMobNotFound          = errors.New("mob not found")
)
