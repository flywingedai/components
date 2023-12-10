package tests

var DefaultMockPriority = -2

/*
Adjust mocks prior to runtime
Default Priority = -2
*/
func (to *TestOptions[C, M, D]) PrepareMocks(
	f func(state *TestState[C, M, D]),
) *TestOptions[C, M, D] {
	testOptions := to.Copy()
	testOptions.options = append(testOptions.options, &testOption[C, M, D]{
		priority: DefaultMockPriority,
		applyFunction: func(state *TestState[C, M, D]) {
			f(state)
		},
	})
	return testOptions
}

/*
Create a mocker object. This object allows normal mocks to be more easily
accessed from test options
*/
func NewMocker() {

}
