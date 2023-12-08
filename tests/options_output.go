package tests

import "github.com/stretchr/testify/assert"

var DefaultOutputPriority = 1

/*
Directly check a value of a particular output.
Default Priority = 1
*/
func (to *TestOptions[C, M, D]) ValidateOutputValue(
	argIndex int,
	expectedValue interface{},
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultOutputPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			assertInterfaceEqual(state.Assertions, expectedValue, state.Output[argIndex])
		},
	})
	return testOptions
}

/*
Specify a value of a particular arg based on the component and/or the test data.
Default Priority = 1
*/
func (to *TestOptions[C, M, D]) ValidateOutput(
	argIndex int,
	f func(state *TestState[C, M, D]) interface{},
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultOutputPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			expectedValue := f(state)
			assertInterfaceEqual(state.Assertions, expectedValue, state.Output[argIndex])
		},
	})
	return testOptions
}

/*
Directly check all values of all outputs at once.
Default Priority = 1
*/
func (to *TestOptions[C, M, D]) ValidateOutputValues(
	expectedValues []interface{},
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultOutputPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			for index, value := range state.Output {
				assertInterfaceEqual(state.Assertions, expectedValues[index], value)
			}
		},
	})
	return testOptions
}

/*
Specify a value of all args at once based on the component and/or the test data.
Default Priority = 1
*/
func (to *TestOptions[C, M, D]) ValidateOutputs(
	f func(state *TestState[C, M, D]) []interface{},
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultOutputPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			expectedValues := f(state)
			for index, value := range state.Output {
				assertInterfaceEqual(state.Assertions, expectedValues[index], value)
			}
		},
	})
	return testOptions
}

/////////////
// HELPERS //
/////////////

// Little helper for ensuring to output values are equal.
func assertInterfaceEqual(parallelAssert *assert.Assertions, expected, actual interface{}) {

	if actual == nil && expected == nil {
		return
	}

	if actual != nil && expected != nil {
		parallelAssert.Equal(expected, actual)
		return
	}

	parallelAssert.Fail("interfaces not equal", actual, expected)

}
