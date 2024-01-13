package tests

import (
	"github.com/stretchr/testify/mock"
)

var DefaultOutputPriority = 10

// Causes an output to be ignore
const Ignore = mock.Anything

// Causes an input or output to be ignored when being set
const Skip = "__tests.skip__"
const Empty = "__tests.empty__"

////////////////
// VALIDATION //
////////////////

/*
After the test has run, check something arbitrary about the test state. If the
provided callback returns an error, the test will fail.
Has Priority = tests.DefaultOutputPriority
*/
func (to *TestOptions[C, M, D]) Validate(
	f func(state *TestState[C, M, D]) error,
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		err := f(state)
		state.Assertions.Nil(err)
	})
}

///////////////////
// SINGLE OUTPUT //
///////////////////

/*
Directly specify the expected value of an output at a given index.
Has Priority = tests.DefaultOutputPriority
Supports DeRef()
*/
func (to *TestOptions[C, M, D]) Output(
	outputIndex int,
	expectedValue interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		assertInterfaceEqual(state.Assertions, handleDereference(expectedValue), state.Output[outputIndex])
	})
}

/*
Specify a pointer to the expected value of an output at a given index.
Has Priority = tests.DefaultOutputPriority
*/
func (to *TestOptions[C, M, D]) Output_P(
	outputIndex int,
	expectedValue *interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		assertInterfaceEqual(state.Assertions, removeInterfacePointer(expectedValue), state.Output[outputIndex])
	})
}

/*
Specify the expected value of an output at a given index based on a provided
function that calculates that value when this option is reached.
Has Priority = tests.DefaultOutputPriority
*/
func (to *TestOptions[C, M, D]) Output_C(
	outputIndex int,
	callbackFunction func() interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		expectedValue := callbackFunction()
		assertInterfaceEqual(state.Assertions, expectedValue, state.Output[outputIndex])
	})
}

/*
Specify the expected value of an output at a given index based on a provided
function that calculates that value based on the value of the TestState when
this option is reached.
Has Priority = tests.DefaultOutputPriority
*/
func (to *TestOptions[C, M, D]) Output_SC(
	outputIndex int,
	callbackFunction func(state *TestState[C, M, D]) interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		expectedValue := callbackFunction(state)
		assertInterfaceEqual(state.Assertions, expectedValue, state.Output[outputIndex])
	})
}

//////////////////////
// MULTIPLE OUTPUTS //
//////////////////////

/*
Directly specify the expected value of all outputs at once.
Has Priority = tests.DefaultOutputPriority
Supports DeRef()
*/
func (to *TestOptions[C, M, D]) Outputs(
	expectedValues ...interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		for index, value := range state.Output {
			if value != Skip {
				assertInterfaceEqual(state.Assertions, expectedValues[index], handleDereference(value))
			}
		}
	})
}

/*
Specify pointers to the expected values of all outputs at once.
Has Priority = tests.DefaultOutputPriority
*/
func (to *TestOptions[C, M, D]) Outputs_P(
	expectedValues ...interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		for index, value := range state.Output {
			if value != Skip {
				assertInterfaceEqual(state.Assertions, removeInterfacePointer(expectedValues[index]), value)
			}
		}
	})
}

/*
Specify the expected values of all outputs at once based on a provided callback
that calculates that value when this option is reached.
Has Priority = tests.DefaultOutputPriority
*/
func (to *TestOptions[C, M, D]) Outputs_C(
	callbackFunction func() []interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		expectedValues := callbackFunction()
		for index, value := range state.Output {
			if value != Skip {
				assertInterfaceEqual(state.Assertions, expectedValues[index], value)
			}
		}
	})
}

/*
Specify the expected values of all outputs at once based on a provided callback
that calculates that value based on the value of the TestState when this option
is reached.
Has Priority = tests.DefaultOutputPriority
*/
func (to *TestOptions[C, M, D]) Outputs_SC(
	callbackFunction func(state *TestState[C, M, D]) []interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultOutputPriority, func(state *TestState[C, M, D]) {
		expectedValues := callbackFunction(state)
		for index, value := range state.Output {
			if value != Skip {
				assertInterfaceEqual(state.Assertions, expectedValues[index], value)
			}
		}
	})
}
