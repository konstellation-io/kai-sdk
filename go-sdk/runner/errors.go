package runner

import (
	"errors"
	"fmt"
)

var (
	ErrLoadingConfiguration = errors.New("configuration could not be loaded %w")
)

// Wrapper creates a function that returns errors starts with a given message.
func Wrapper(message string) func(params ...interface{}) error {
	return func(params ...interface{}) error {
		return fmt.Errorf(message, params...) //nolint: goerr113
	}
}
