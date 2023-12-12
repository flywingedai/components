package tests

var DefaultMockPriority = -2

/*
Adjust mocks prior to runtime
Default Priority = -2
*/
func (to *TestOptions[C, M, D]) PrepareMocks(
	f func(state *TestState[C, M, D]),
) *TestOptions[C, M, D] {
	return to.copyAndAppend(DefaultMockPriority, func(state *TestState[C, M, D]) {
		f(state)
	})
}
