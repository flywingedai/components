package tests

import (
	"reflect"
	"sort"

	"github.com/stretchr/testify/assert"
)

/*
Individual test to be run by a tester. This is the interface which all
children test types must conform to.
*/
type TestConfig[C, M, D any] struct {
	name            string
	getTestFunction func(state *TestState[C, M, D]) reflect.Value
	Options         *TestOptions[C, M, D]
}

/*
The internal test state. Is automatically updated as options are called.
*/
type TestState[C, M, D any] struct {
	Assertions *assert.Assertions

	Component C
	Mocks     *M
	Data      *D

	Input  []interface{}
	Output []interface{}
}

/*
Helper function for adding a test to a tester. Makes it easier to do this inline
*/
func (tc *TestConfig[C, M, D]) Register(tester *Tester[C, M, D]) *TestConfig[C, M, D] {
	tester.AddTests(tc)
	return tc
}

/*
Standard run command for all states
*/
func (tc *TestConfig[C, M, D]) run(state *TestState[C, M, D]) {

	// Order all the options in ascending priority
	sort.SliceStable(tc.Options.options, func(i, j int) bool {
		return tc.Options.options[i].priority < tc.Options.options[j].priority
	})

	// Loop through all the options with value < 0
	runFunctionIndex := 0
	for i, option := range tc.Options.options {

		// We stop processing at this point as we now need to
		if option.priority >= 0 {
			runFunctionIndex = i
			break
		}

		// Otherwise, we just run the applyFunction for this option
		option.applyFunction(state)

	}

	// Now we fetch the test function and run it with the args
	f := tc.getTestFunction(state)
	args := getCallArgs(state.Input)
	reflectOutput := f.Call(args)
	state.Output = make([]interface{}, len(reflectOutput))

	// Convert the reflected output to normal interface output
	for i := range reflectOutput {
		state.Output[i] = getReflectInterface(reflectOutput[i])
	}

	// Loop through all the options with value >= 0
	for _, option := range tc.Options.options[runFunctionIndex:] {
		option.applyFunction(state)
	}

}

/*
Helper function that converts the input data array into something the reflect
package knows how to parse with the Call() method
*/
func getCallArgs(input []interface{}) []reflect.Value {

	// Convert the input into reflect types so it can be fed into the real function
	args := []reflect.Value{}
	for _, arg := range input {

		if arg == nil {

			/*
				If an arg is nil, we just need a nil pointer of any
				type. This will allow the reflect package to
				correctly make the call. I can't think of any cases
				where this doesn't work given that when doing type
				checking of a nil value you will get an error anyway...
			*/
			var nilSubstitute *int
			args = append(args, reflect.ValueOf(nilSubstitute))

		} else {

			/*
				If we don't have a nil value, we just grab the
				reflect.ValueOf the arg and add it to the args array.
			*/
			args = append(args, reflect.ValueOf(arg))
		}

	}

	return args

}

// Helper to get a meaningful interface out of a reflect value
func getReflectInterface(value reflect.Value) interface{} {
	interfaceValue := value.Interface()

	// Because of the reflect package, we have to do some fancy
	// jiggery-pokery to avoid bad comparisons
	isNil := false
	if interfaceValue == nil {
		return nil
	}

	switch reflect.TypeOf(interfaceValue).Kind() {
	case reflect.Array, reflect.Map, reflect.Ptr, reflect.Chan, reflect.Slice:
		isNil = reflect.ValueOf(interfaceValue).IsNil()
	}

	if isNil {
		return nil
	} else {
		return interfaceValue
	}

}
