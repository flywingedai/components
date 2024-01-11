package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var DefaultOutputPriority = 1

// Causes an output to be ignore
const IGNORE = mock.Anything

// Causes an input or output to be ignored when being set
const SKIP = "__tests.skipInput__"

func (to *TestOptions[C, M, D]) Validate(
	f func(state *TestState[C, M, D]) error,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		err := f(state)
		state.Assertions.Nil(err)
	})
}

/*
Directly check a value of a particular output.
Default Priority = 1
*/
func (to *TestOptions[C, M, D]) ValidateOutputValue(
	argIndex int,
	expectedValue interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		assertInterfaceEqual(state.Assertions, expectedValue, state.Output[argIndex])
	})
}

/*
Directly check a value of a particular output.
Default Priority = 1
*/
func (to *TestOptions[C, M, D]) ValidateOutputValue_Pointer(
	argIndex int,
	expectedValue *interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		assertInterfaceEqual(state.Assertions, removeInterfacePointer(expectedValue), state.Output[argIndex])
	})
}

/*
Specify a value of a particular arg based on the component and/or the test data.
Default Priority = 1
*/
func (to *TestOptions[C, M, D]) ValidateOutput(
	argIndex int,
	f func(state *TestState[C, M, D]) interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		expectedValue := f(state)
		assertInterfaceEqual(state.Assertions, expectedValue, state.Output[argIndex])
	})
}

/*
Directly check all values of all outputs at once.
Default Priority = 1
*/
func (to *TestOptions[C, M, D]) ValidateOutputValues(
	expectedValues ...interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		for index, value := range state.Output {
			if value != SKIP {
				assertInterfaceEqual(state.Assertions, expectedValues[index], value)
			}
		}
	})
}

/*
Directly check all values of all outputs at once.
Default Priority = 1
*/
func (to *TestOptions[C, M, D]) ValidateOutputValues_Pointer(
	expectedValues ...interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		for index, value := range state.Output {
			if value != SKIP {
				assertInterfaceEqual(state.Assertions, removeInterfacePointer(expectedValues[index]), value)
			}
		}
	})
}

/*
Specify a value of all args at once based on the component and/or the test data.
Default Priority = 1
*/
func (to *TestOptions[C, M, D]) ValidateOutputs(
	f func(state *TestState[C, M, D]) []interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		expectedValues := f(state)
		for index, value := range state.Output {
			if value != SKIP {
				assertInterfaceEqual(state.Assertions, expectedValues[index], value)
			}
		}
	})
}

/////////////
// HELPERS //
/////////////

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
