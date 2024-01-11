package tests

////////////
// REPEAT //
////////////

// Repeat all the options in this list "count" times
func (to *TestOptions[C, M, D]) Repeat(
	count int,
) *TestOptions[C, M, D] {
	newOptions := to.Copy()
	for i := 0; i < count; i++ {
		newOptions.options = append(newOptions.options, newOptions.options...)
	}
	return to
}

// Repeat the last option added in "count" times
func (to *TestOptions[C, M, D]) RepeatLast(
	count int,
) *TestOptions[C, M, D] {
	newOptions := to.Copy()
	copiedOption := to.options[len(to.options)-1]
	for i := 0; i < count; i++ {
		newOptions.options = append(newOptions.options, copiedOption)
	}
	return to
}

// Repeat the last "n" options added in "count" times
func (to *TestOptions[C, M, D]) RepeatLastN(
	count int,
	n int,
) *TestOptions[C, M, D] {
	newOptions := to.Copy()
	copiedOptions := to.options[len(to.options)-n:]
	for i := 0; i < count; i++ {
		newOptions.options = append(newOptions.options, copiedOptions...)
	}
	return to
}
