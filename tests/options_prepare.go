package tests

var DefaultPreparePriority = -30
var DefaultSetupPriority = -50

/*
Perform some generic action.
Has Priority = tests.DefaultPreparePriority
*/
func (to *TestOptions[C, M, D]) Prepare(
	prepareFunction func(),
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultPreparePriority, func(_ *TestState[C, M, D]) {
		prepareFunction()
	})
}

/*
Prepare some action dependant on the data in the internal TestState.
Has Priority = tests.DefaultPreparePriority
*/
func (to *TestOptions[C, M, D]) Prepare_D(
	prepareFunction func(data *D),
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultPreparePriority, func(state *TestState[C, M, D]) {
		prepareFunction(state.Data)
	})
}

/*
Prepare some action dependant on the entirety of the internal TestState.
Has Priority = tests.DefaultPreparePriority
*/
func (to *TestOptions[C, M, D]) Prepare_S(
	prepareFunction func(state *TestState[C, M, D]),
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultPreparePriority, func(state *TestState[C, M, D]) {
		prepareFunction(state)
	})
}
