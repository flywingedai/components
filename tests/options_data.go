package tests

var DefaultDataPriority = -3

/*
Adjust mocks prior to runtime without care for test data or component state
Default Priority = -3
*/
func (to *TestOptions[C, M, D]) PrepareData(
	f func(data *D),
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultDataPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			f(state.Data)
		},
	})
	return testOptions
}
