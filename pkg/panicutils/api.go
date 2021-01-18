package panicutils

import (
	"fmt"
)

// PanicIfError will panic if the provided error is not nil
func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

// RecoverToErrorAndHandle recovers a panic, converts it to an error and invokes the handler
func RecoverToErrorAndHandle(handler func(err error)) {
	if err := PanicToError(recover()); err != nil {
		handler(err)
	}
}

// PanicToError converts a panic to an error
func PanicToError(recovered interface{}) error {
	if recovered != nil {
		if err, ok := recovered.(error); ok {
			return err
		}
		return fmt.Errorf("panic: %v", recovered)
	}
	return nil
}
