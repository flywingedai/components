package tests

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
A tester is the base unit that handles mock testing for a component.
Type P is the component params
Type C is the component type
Type M is the mocks type
Type D is the data type
*/
type Tester[P, C, M, D any] struct {

	/*
		Internal testing state. Contains all the different tests and the options
		that are associated with each one.
	*/
	tests []*testConfig[C, M, D]

	// Params passed in during Tester creation
	newComponentFunction func(P) C
	buildMocksFunction   func(*testing.T) (P, *M)
	initDataFunction     func() *D

	// Global options for this tester
	Options *TestOptions[C, M, D]
}

/*
Create a new Tester with a specified Component, Mocks, and Data structure. Use
tests.NullDataInitialization as the initDataFunction if you do not need
to use the provided data object to facilitate your tester.
*/
func NewTester[P, C, M, D any](
	newComponentFunction func(P) C,
	buildMocksFunction func(*testing.T) (P, *M),
	initDataFunction func() *D,
) *Tester[P, C, M, D] {
	tester := &Tester[P, C, M, D]{
		newComponentFunction: newComponentFunction,
		buildMocksFunction:   buildMocksFunction,
		initDataFunction:     initDataFunction,
		Options:              &TestOptions[C, M, D]{},
	}

	return tester
}

/*
Create a new Tester with a specified Component, Mocks, and Data structure. Use
tests.NullDataInitialization as the initDataFunction if you do not need
to use the provided data object to facilitate your tester.
*/
func NewGinTester[P, C, M, D any](
	newComponentFunction func(P) C,
	buildMocksFunction func(*testing.T) (P, *M),
	initDataFunction func() *D,
) *Tester[P, C, M, D] {
	tester := &Tester[P, C, M, D]{
		newComponentFunction: newComponentFunction,
		buildMocksFunction:   buildMocksFunction,
		initDataFunction:     initDataFunction,
		Options:              &TestOptions[C, M, D]{},
	}

	tester.Options = tester.Options.SetInput(0, func(state *TestState[C, M, D]) interface{} {
		return convertToGinDataInterface(state.Data).GetCtx()
	})

	return tester
}

/*
Create a new Tester with a specified Component and Mocks structure.
No initialization step is called at this point.
*/
func NewTesterWithoutData[P, C, M any](
	newComponentFunction func(P) C,
	buildMocksFunction func(*testing.T) (P, *M),
) *Tester[P, C, M, interface{}] {
	tester := &Tester[P, C, M, interface{}]{
		newComponentFunction: newComponentFunction,
		buildMocksFunction:   buildMocksFunction,
		initDataFunction: func() *interface{} {
			return nil
		},
	}

	return tester
}

/*
Create a new test that runs a method of the parent component
*/
func (tester *Tester[P, C, M, D]) NewMethodTest(
	testName, methodName string,
	options *TestOptions[C, M, D],
) {
	tester.tests = append(tester.tests, &testConfig[C, M, D]{
		name: testName,
		getTestFunction: func(c C) reflect.Value {
			return reflect.ValueOf(c).MethodByName(methodName)
		},
		options: options,
	})
}

/*
Create a new test that runs an arbitrary, non-method function
*/
func (tester *Tester[P, C, M, D]) NewFunctionTest(
	testName string,
	function interface{},
	options []*testOption[C, M, D],
) {
	tester.tests = append(tester.tests, &testConfig[C, M, D]{
		name: testName,
		getTestFunction: func(_ C) reflect.Value {
			return reflect.ValueOf(function)
		},
		options: &TestOptions[C, M, D]{
			options: options,
		},
	})
}

/*
Runs all the currently generated tests
*/
func (tester *Tester[P, C, M, D]) Test(t *testing.T) {
	for _, loopTest := range tester.tests {
		test := loopTest

		// Create the test to run
		t.Run(test.name, func(t *testing.T) {
			fmt.Println(test.name)

			// Support parallel tests running at the same time
			t.Parallel()
			a := assert.New(t)

			// Create the base objects for the test
			c, m, d := tester.build(t)
			test.options = tester.Options.Combine(test.options)

			// Run the test with the parallel assertion and base objects
			testState := TestState[C, M, D]{
				Assertions: a,

				Component: c,
				Mocks:     m,
				Data:      d,

				Input:  []interface{}{},
				Output: []interface{}{},
			}
			test.run(&testState)
		})
	}
}

/*
Creates a component and a mocks object to user for the testing
*/
func (tester *Tester[P, C, M, D]) build(t *testing.T) (C, *M, *D) {
	params, mocks := tester.buildMocksFunction(t)
	component := tester.newComponentFunction(params)
	data := tester.initDataFunction()
	return component, mocks, data
}
