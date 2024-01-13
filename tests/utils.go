package tests

import (
	"errors"
	"reflect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

// Function to convert an array from one type to another
func mapArray[I, O any](input []I, convertFunction func(I) O) []O {
	output := []O{}
	for _, value := range input {
		output = append(output, convertFunction(value))
	}
	return output
}

// Function to convert a map from one type to another
func mapMap[K comparable, VI, VO any](input map[K]VI, convertFunction func(VI) VO) map[K]VO {
	output := map[K]VO{}
	for key, value := range input {
		output[key] = convertFunction(value)
	}
	return output
}

// Little helper for ensuring to output values are equal.
func assertInterfaceEqual(parallelAssert *assert.Assertions, expected, actual interface{}) {

	if expected == mock.Anything {
		return
	}

	if actual == nil && expected == nil {
		return
	}

	if actual != nil && expected != nil {
		parallelAssert.Equal(expected, actual)
		return
	}

	parallelAssert.Fail("interfaces not equal", actual, expected)

}
