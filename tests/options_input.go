package tests

var DefaultInputPriority = -1

/*
Directly specify a value of a particular arg.
Default Priority = -1
*/
func (to *TestOptions[C, M, D]) SetInputByValue(
	argIndex int,
	value interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		state.Input = expandInput(state.Input, argIndex)
		(state.Input)[argIndex] = value
	})
}

/*
Directly specify a value of a particular arg.
Default Priority = -1
*/
func (to *TestOptions[C, M, D]) SetInputByValue_Pointer(
	argIndex int,
	value interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		state.Input = expandInput(state.Input, argIndex)
		(state.Input)[argIndex] = removeInterfacePointer(value)
	})
}

/*
Specify a value of a particular arg based on the component and/or the test data.
Default Priority = -1
*/
func (to *TestOptions[C, M, D]) SetInput(
	argIndex int,
	f func(state *TestState[C, M, D]) interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		state.Input = expandInput(state.Input, argIndex)
		(state.Input)[argIndex] = f(state)
	})
}

/*
Directly specify a value of all args at once.
Default Priority = -1
*/
func (to *TestOptions[C, M, D]) SetInputsByValue(
	values ...interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		state.Input = make([]interface{}, len(values))
		for i, value := range values {
			if value != SKIP {
				state.Input[i] = value
			}
		}
	})
}

/*
Directly specify a value of all args at once.
Default Priority = -1
*/
func (to *TestOptions[C, M, D]) SetInputsByValue_Pointer(
	values ...interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		state.Input = make([]interface{}, len(values))
		for i, value := range values {
			if value != SKIP {
				state.Input[i] = removeInterfacePointer(value)
			}
		}
	})
}

/*
Specify a value of all args at once based on the component and/or the test data.
Default Priority = -1
*/
func (to *TestOptions[C, M, D]) SetInputs(
	f func(state *TestState[C, M, D]) []interface{},
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultInputPriority, func(state *TestState[C, M, D]) {
		values := f(state)
		state.Input = make([]interface{}, len(values))
		for i, value := range values {
			if value != SKIP {
				state.Input[i] = value
			}
		}
	})
}

/////////////
// HELPERS //
/////////////

// Little helper for managing inputs
func expandInput(input []interface{}, size int) []interface{} {
	newInput := make([]interface{}, size)
	for i := 0; i < size; i++ {
		if i < len(input) {
			newInput[i] = input[i]
		} else {
			newInput[i] = SKIP
		}
	}
	return newInput
}
