package tests

import (
	"errors"
	"reflect"
)

// Helper to handle bad read errors
type ErrorIO struct {
	ReadErr  error
	CloseErr error
}

func NewDefaultErrorIO() *ErrorIO {
	return &ErrorIO{
		ReadErr:  errors.New("ErrorIO: Read error"),
		CloseErr: errors.New("ErrorIO: Close error"),
	}
}

func (e *ErrorIO) Read(p []byte) (n int, err error) {
	return 0, e.ReadErr
}

func (e *ErrorIO) Close() error {
	return e.CloseErr
}

func removeInterfacePointer(pointerInterface interface{}) interface{} {

	if pointerInterface == nil {
		return nil
	}

	reflectValue := reflect.ValueOf(pointerInterface)

	if reflectValue.Kind() == reflect.Pointer {
		elem := reflectValue.Elem()
		return elem.Interface()
	}

	return pointerInterface
}

/*
Exposed interface pointer removal that requires the type to be passed into a
generic function so a default value can be created in the zero case.
*/
func RemoveInterfacePointer[T any](pointerInterface interface{}) interface{} {

	if pointerInterface == nil {
		return nil
	}

	reflectValue := reflect.ValueOf(pointerInterface)

	if reflectValue.Kind() == reflect.Pointer {
		elem := reflectValue.Elem()

		if elem.IsZero() {
			var zeroValue T
			return zeroValue
		}

		return elem.Interface()
	}

	return pointerInterface
}
