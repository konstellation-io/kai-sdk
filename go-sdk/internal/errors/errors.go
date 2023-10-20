package errors

import (
	"errors"
	"fmt"
)

var (
	ErrUndefinedEphemeralStorage = errors.New("the ephemeral storage does not exist")
	ErrMessageToBig              = errors.New("compressed message exceeds maximum size allowed")
	ErrMsgAck                    = "Error in message ack" //nolint:gochecknoglobals // This is a constant
	ErrEmptyPayload              = errors.New("the payload cannot be empty")
	ErrObjectAlreadyExists       = errors.New("object already exists for the given key")
)

// Wrapper creates a function that returns errors starts with a given message.
func Wrapper(message string) func(params ...interface{}) error {
	return func(params ...interface{}) error {
		return fmt.Errorf(message, params...) //nolint:goerr113 // This is a wrapper
	}
}
