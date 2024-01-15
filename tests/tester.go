package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
A tester is the base unit that handles mock testing for a component.
Type C is the base component type
Type M is the component mocks type
Type D is the data type.
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

	// List of branches associated with this tester
	branches map[string]*TestOptions[C, M, D]

	// Global options for this tester
	Options *TestOptions[C, M, D]

	/*
		Whether or not the tester should run tests in parallel. This is not
		recommended unless you are utilizing test data.
	*/
	Parallel bool
}

/////////////////
// NEW TESTERS //
/////////////////

// Internal helper function for making a new tester
func emptyTester[C, M, D any]() *Tester[C, M, D] {
	tester := &Tester[C, M, D]{
		Options:  &TestOptions[C, M, D]{},
		branches: map[string]*TestOptions[C, M, D]{},
	}
	tester.Options.tester = tester
	return tester
}

/*
Create a new Tester with a specified Component, Mocks, and Data structure.
*/
func NewTesterWithData[C, M, D any](
	buildMocksFunction func(*testing.T) (C, *M),
	initDataFunction func() *D,
) *Tester[C, M, D] {
	tester := emptyTester[C, M, D]()
	tester.buildMocksFunction = buildMocksFunction
	tester.initDataFunction = initDataFunction
	return tester
}

/*
Create a new Tester with a specified Component and Mocks structure. Requires
a provided initialization function.
*/
func NewTesterWithInit[C, M any](
	buildMocksFunction func(*testing.T) (C, *M),
	initDataFunction func(),
) *Tester[C, M, interface{}] {
	tester := emptyTester[C, M, interface{}]()
	tester.buildMocksFunction = buildMocksFunction
	tester.initDataFunction = func() *interface{} {
		initDataFunction()
		return nil
	}
	return tester
}

/*
Create a new Tester with a specified Component and Mocks structure.
No initialization step is called for this kind of tester. The inferred
type for the data is interface{}, but it will always be set to nil for
tests created this way.
*/
func NewTesterWithoutInit[C, M any](
	buildMocksFunction func(*testing.T) (C, *M),
) *Tester[C, M, interface{}] {
	tester := emptyTester[C, M, interface{}]()
	tester.buildMocksFunction = buildMocksFunction
	tester.initDataFunction = func() *interface{} { return nil }
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
	tester := emptyTester[interface{}, interface{}, D]()
	tester.buildMocksFunction = func(t *testing.T) (interface{}, *interface{}) { return nil, nil }
	tester.initDataFunction = initDataFunction
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
func NewFunctionTesterWithoutData() *Tester[interface{}, interface{}, interface{}] {
	tester := emptyTester[interface{}, interface{}, interface{}]()
	tester.buildMocksFunction = func(t *testing.T) (interface{}, *interface{}) { return nil, nil }
	tester.initDataFunction = func() *interface{} { return nil }
	return tester
}

////////////////////
// TESTER METHODS //
////////////////////

/*
Create a new options for this tester without any of the existing options
included. Makes it slightly easier to create branches.
*/
func (tester *Tester[C, M, D]) NewOptions() *TestOptions[C, M, D] {
	return NewOptions[C, M, D](tester)
}

/*
Checkout the tagged TestOptions branch.

Once a tag is applied, you can fetch it:
from any child TestOptions by: testOptions.Checkout($tagName)
or from the parent tester by: tester.Checkout($tagName)
*/
func (tester *Tester[C, M, D]) Checkout(
	tag string,
) *TestOptions[C, M, D] {
	options, exists := tester.branches[tag]
	if !exists {
		panic("could not find tag " + tag + " in tester")
	}
	return options
}

/*
Register tests with the tester. All registered tests will be run when
tester.Test(t) is called
*/
func (tester *Tester[C, M, D]) RegisterTests(
	tests ...*TestConfig[C, M, D],
) *Tester[C, M, D] {
	tester.tests = append(tester.tests, tests...)
	return tester
}

/*
Attach a group id to the tester so that all test names under this tester
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
			if tester.Parallel {
				t.Parallel()
			}
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
