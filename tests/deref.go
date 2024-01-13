package tests

/*
Internal object that is produced by the tests.Deref() function. Used to
indicate to the tests.handleDerefernece() function whether or not this value
should be dereferenced or not.
*/
type deref struct {
	dereference func() interface{}
}

/*
Deref should be used inside of functions that have "Supports DeRef()" written
in their docstring somewhere. It allows you to pass in a pointer to something
that you would like to be evaluated during the test, not before.
*/
func DeRef(field interface{}) deref {
	return deref{
		dereference: func() interface{} {
			return removeInterfacePointer(field)
		},
	}
}

/*
Internal function for managing dereferences of interface values. If the value
is of type tests.deref, then it will be dereferenced, otherwise it will be
returned as is.
*/
func handleDereference(value interface{}) interface{} {
	derefStruct, ok := value.(deref)
	if ok {
		return derefStruct.dereference()
	}
	return value
}
