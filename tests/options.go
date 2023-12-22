package tests

import "reflect"

/*
The test option is the main struct for handling test options. The applyFunction
should modify the args that are passed in, rather than returning anything.
*/
type testOption[C, M, D any] struct {

	/*
		When this option's applyFunction will be called relative to other
		options that happen at the same time. Negative values occur before
		the test function is run, positive values occur afterwards.
		Zero values are reserved for options which set the input of the test
		function.
	*/
	priority int

	// Used to facilitate option checkouts for easier iterative test creation.
	tag string

	/*
		The applyFunction is what gets called by the TestOption interface to
		make sure the test ran nominally
	*/
	applyFunction func(state *TestState[C, M, D])
}

//////////////////
// TEST OPTIONS //
//////////////////

/*
The TestOptions object is responsible providing binding to all the useful test
options available. It implemented method chaining in order to make for a more
simple interface for creating options.

Each chained option should return a completely new TestOptions struct so as to
allow easy bifrucation of test setup.
*/
type TestOptions[C, M, D any] struct {
	options []*testOption[C, M, D]
}

func NewOptions[C, M, D any]() *TestOptions[C, M, D] {
	return &TestOptions[C, M, D]{
		options: []*testOption[C, M, D]{},
	}
}

// Helper copy function
func (to *TestOptions[C, M, D]) Copy() *TestOptions[C, M, D] {
	testOptions := []*testOption[C, M, D]{}
	testOptions = append(testOptions, to.options...)
	return &TestOptions[C, M, D]{
		options: testOptions,
	}
}

/*
Append creates a new TestOptions object with all the options provided combined.
*/
func (to *TestOptions[C, M, D]) Append(
	otherTestOptions ...*TestOptions[C, M, D],
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	for _, option := range otherTestOptions {
		testOptions.options = append(testOptions.options, option.options...)
	}
	return testOptions
}

/*
Adjusts the priority of the last option in the TestOptions. If no options exist
currently, will panic so that the method chaining archetype can be preserved.
Priorities <  0 are run before the test runs
Priorities > 0 are run after test runs
Priority = 0 is reserved for state cleanup before checking outputs
Lower value priority occurs first.
*/
func (to *TestOptions[C, M, D]) WithPriority(priority int) *TestOptions[C, M, D] {
	if len(to.options) > 0 {
		to.options[len(to.options)-1].priority = priority
	} else {
		panic("called TestOptions.WithPriority() with no options.")
	}
	return to
}

/*
General use case for adding an option. If there is any arbitrary thing you need
to do, and the premade functions do not allow you to do it, you can add a
generic option to allow you to do anything. It sends the whole state of the test
into the "applyFunction" for you to modify as you please.
*/
func (to *TestOptions[C, M, D]) NewOption(
	priority int,
	applyFunction func(state *TestState[C, M, D]),
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority:      priority,
		applyFunction: applyFunction,
	})
	return testOptions
}

/*
Create a tag which allows options to be "checked-out" later. This
allows more simple control over options when writing iterative
tests on each other.
*/
func (to *TestOptions[C, M, D]) Tag(
	tag string,
) *TestOptions[C, M, D] {
	to.options[len(to.options)-1].tag = tag
	return to
}

/*
Checkout the most recent entry of the tag in the options. Return
all options prior to that tag, including that tag. If the tag is
not found, this will panic.
*/
func (to *TestOptions[C, M, D]) Checkout(
	tag string,
) *TestOptions[C, M, D] {
	for i := len(to.options) - 1; i >= 0; i-- {
		if to.options[i].tag == tag {
			return &TestOptions[C, M, D]{
				options: to.options[:i+1],
			}
		}
	}
	panic("could not find tag " + tag + " in TestOptions")
}

func (to *TestOptions[C, M, D]) copyAndAppend(
	priority int,
	applyFunction func(state *TestState[C, M, D]),
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority:      priority,
		applyFunction: applyFunction,
	})
	return testOptions
}

/*
Create a new test based on the given function. By doing things this way, private
methods can be tested as well by accessing them from the state at runtime.
*/
func (to *TestOptions[C, M, D]) CreateTest(
	testName string,
	getTestFunction func(state *TestState[C, M, D]) interface{},
) *TestConfig[C, M, D] {
	return &TestConfig[C, M, D]{
		name: testName,
		getTestFunction: func(state *TestState[C, M, D]) reflect.Value {
			return reflect.ValueOf(getTestFunction(state))
		},
		Options: to,
	}
}

/*
Create a new test which automatically fetches the named component
method at runtime.
*/
func (to *TestOptions[C, M, D]) CreateMethodTest(method, testName string) *TestConfig[C, M, D] {
	return &TestConfig[C, M, D]{
		name: testName,
		getTestFunction: func(state *TestState[C, M, D]) reflect.Value {
			return reflect.ValueOf(state.Component).MethodByName(method)
		},
		Options: to,
	}
}

/*
Create a new test which automatically fetches the function given
at runtime.
*/
func (to *TestOptions[C, M, D]) CreateFunctionTest(function interface{}, testName string) *TestConfig[C, M, D] {
	return &TestConfig[C, M, D]{
		name: testName,
		getTestFunction: func(_ *TestState[C, M, D]) reflect.Value {
			return reflect.ValueOf(function)
		},
		Options: to,
	}
}
