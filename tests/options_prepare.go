package tests

var DefaultPreparePriority = -3

/*
Adjust data prior to runtime without care for test state
Default Priority = -3
*/
func (to *TestOptions[C, M, D]) Prepare(
	f func(),
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultPreparePriority, func(state *TestState[C, M, D]) {
		f()
	})
}

/*
Adjust data prior to runtime without care for mocks or component state
Default Priority = -3
*/
func (to *TestOptions[C, M, D]) PrepareData(
	f func(data *D),
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultPreparePriority, func(state *TestState[C, M, D]) {
		f(state.Data)
	})
}

/*
Adjust state prior to other options
Default Priority = -3
*/
func (to *TestOptions[C, M, D]) PrepareState(
	f func(state *TestState[C, M, D]),
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultPreparePriority, func(state *TestState[C, M, D]) {
		f(state)
	})
}
