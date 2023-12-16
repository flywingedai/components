package tests

import (
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
type Tester[C, M, D any] struct {

	/*
		Internal testing state. Contains all the different tests and the options
		that are associated with each one.
	*/
	testGroupID string
	tests       []*TestConfig[C, M, D]

	// Params passed in during Tester creation
	buildMocksFunction func(*testing.T) (C, *M)
	initDataFunction   func() *D

	// Global options for this tester
	Options *TestOptions[C, M, D]
}

/*
Create a new Tester with a specified Component, Mocks, and Data structure. Use
tests.NullDataInitialization as the initDataFunction if you do not need
to use the provided data object to facilitate your tester.
*/
func NewTester[C, M, D any](
	buildMocksFunction func(*testing.T) (C, *M),
	initDataFunction func() *D,
) *Tester[C, M, D] {
	tester := &Tester[C, M, D]{
		buildMocksFunction: buildMocksFunction,
		initDataFunction:   initDataFunction,
		Options:            &TestOptions[C, M, D]{},
	}

	return tester
}

/*
Create a new Tester for a function or a group of functions. This
does not require any components or mocks to run as it is intended
to be used on individual functions. The rest of the tester suite
can still be used in this case, although the inferred generic types
for the C (Component) and M (Mocks) fields will both simply be
interfaces. You should ignore those fields in the Options.
*/
func NewFunctionTester[D any](
	initDataFunction func() *D,
) *Tester[interface{}, interface{}, D] {
	tester := &Tester[interface{}, interface{}, D]{
		buildMocksFunction: func(t *testing.T) (interface{}, *interface{}) { return nil, nil },
		initDataFunction:   initDataFunction,
		Options:            &TestOptions[interface{}, interface{}, D]{},
	}

	return tester
}

/*
Create a new Tester with a specified Component and Mocks structure.
No initialization step is called for this kind of tester. The inferred
type for the data is interface{}, but it will always be set to nil for
tests created this way.
*/
func NewTesterWithoutData[C, M any](
	buildMocksFunction func(*testing.T) (C, *M),
) *Tester[C, M, interface{}] {
	tester := &Tester[C, M, interface{}]{
		buildMocksFunction: buildMocksFunction,
		initDataFunction: func() *interface{} {
			return nil
		},
		Options: &TestOptions[C, M, interface{}]{},
	}

	return tester
}

/*
Add tests that run a method of the parent component
*/
func (tester *Tester[C, M, D]) AddTests(
	tests ...*TestConfig[C, M, D],
) *Tester[C, M, D] {
	tester.tests = append(tester.tests, tests...)
	return tester
}

/*
Attach a group id to the tester so that all tests under this tester
automatically have a prefix attached to them.
*/
func (tester *Tester[C, M, D]) WithGroupID(
	groupID string,
) *Tester[C, M, D] {
	tester.testGroupID = groupID
	return tester
}

/*
Runs all the currently appended tests in the order in which they were appended.
*/
func (tester *Tester[C, M, D]) Test(t *testing.T) {
	for _, loopTest := range tester.tests {
		test := loopTest

		// Determine if there should be a prefix for the test
		name := test.name
		if tester.testGroupID != "" {
			name = tester.testGroupID + ": " + name
		}

		// Create the test to run
		t.Run(name, func(t *testing.T) {

			// Support parallel tests running at the same time
			t.Parallel()
			a := assert.New(t)

			// Create the base objects for the test
			c, m, d := tester.build(t)
			test.Options = tester.Options.Append(test.Options)

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
Creates a component and a mocks object for testing. Is automatically called at
the beginning of each test so that a fresh set of components, mocks, and data
is available.
*/
func (tester *Tester[C, M, D]) build(t *testing.T) (C, *M, *D) {
	component, mocks := tester.buildMocksFunction(t)
	data := tester.initDataFunction()
	return component, mocks, data
}
