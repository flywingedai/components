package tests

var DefaultPreparePriority = -3

/*
Adjust data prior to runtime without care for test data or component state
Default Priority = -3
*/
func (to *TestOptions[C, M, D]) PrepareData(
	f func(data *D),
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultPreparePriority,
		applyFunction: func(state *TestState[C, M, D]) {
			f(state.Data)
		},
	})
	return testOptions
}

/*
Adjust state prior to other options
Default Priority = -3
*/
func (to *TestOptions[C, M, D]) PrepareState(
	f func(state *TestState[C, M, D]),
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultPreparePriority,
		applyFunction: func(state *TestState[C, M, D]) {
			f(state)
		},
	})
	return testOptions
}
