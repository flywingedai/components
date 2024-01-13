package tests

var DefaultInputPriority = -10

//////////////////
// SINGLE INPUT //
//////////////////

// Internal function to set a single input value
func applySetInput[C, M, D any](state *TestState[C, M, D], argIndex int, value interface{}) {
	state.Input = expandInput(state.Input, argIndex)
	state.Input[argIndex] = value
}

/*
Directly specify the value of a particular arg.
Has Priority = tests.DefaultInputPriority
Supports DeRef()
*/
func (to *TestOptions[C, M, D]) SetInput(
	argIndex int,
	value interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applySetInput(state, argIndex, handleDereference(value))
	})
}

/*
Specify the pointer to the value of a particular arg. If the input expects the
type "int", you would pass in "*int".
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) SetInput_P(
	argIndex int,
	value interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applySetInput(state, argIndex, removeInterfacePointer(value))
	})
}

/*
Specify a value of a particular arg based on a provided callback that calculates
that value when this option is reached.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) SetInput_C(
	argIndex int,
	callbackFunction func() interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applySetInput(state, argIndex, callbackFunction())
	})
}

/*
Specify a value of a particular arg based on a provided callback that calculates
that value based on the value of the TestState when this option is reached.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) SetInput_SC(
	argIndex int,
	callbackFunction func(state *TestState[C, M, D]) interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applySetInput(state, argIndex, callbackFunction(state))
	})
}

/////////////////////
// MULTIPLE INPUTS //
/////////////////////

// Internal function to set multiple input values at once
func applySetInputs[C, M, D any](state *TestState[C, M, D], values []interface{}) {
	state.Input = expandInput(state.Input, len(values)-1)
	for i, value := range values {
		if value != Skip {
			state.Input[i] = value
		}
	}
}

/*
Directly specify the value of all args at once. To skip setting a particular
arg, set it to tests.Skip
Has Priority = tests.DefaultInputPriority
Supports DeRef()
*/
func (to *TestOptions[C, M, D]) SetInputs(
	values ...interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applySetInputs(state, mapArray(values, handleDereference))
	})
}

/*
Specify the pointer to all arg values at once. To skip setting a particular
arg, set it to tests.Skip
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) SetInputs_P(
	values ...interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applySetInputs(state, mapArray(values, removeInterfacePointer))
	})
}

/*
Specify all args valuesbased on a provided callback that calculates those values
when this option is reached.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) SetInputs_C(
	callbackFunction func() []interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applySetInputs(state, callbackFunction())
	})
}

/*
Specify all args valuesbased on a provided callback that calculates those values
based on the value of the TestState when this option is reached.
Has Priority = tests.DefaultInputPriority
*/
func (to *TestOptions[C, M, D]) SetInputs_SC(
	callbackFunction func(state *TestState[C, M, D]) []interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		applySetInputs(state, callbackFunction(state))
	})
}

/////////////
// HELPERS //
/////////////

// Little helper for managing the size of TestState.Input
func expandInput(input []interface{}, argIndex int) []interface{} {

	/*
		Determine how large the resulting input needs to be based on whether or
		not the argIndex or the existing size is larger.
	*/
	newInputSize := argIndex + 1
	if len(input) > newInputSize {
		newInputSize = len(input)
	}

	/*
		Create the new input by copying over old data. Any newly created data
		will be set to tests.Empty.
	*/
	newInput := make([]interface{}, newInputSize)
	for i := 0; i < newInputSize; i++ {
		if i < len(input) {
			newInput[i] = input[i]
		} else {
			newInput[i] = Empty
		}
	}
	return newInput
}
