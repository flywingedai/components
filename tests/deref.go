package tests

/////////////////
// DEREFERENCE //
/////////////////

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

func handleDereference(value interface{}) interface{} {
	derefStruct, ok := value.(deref)
	if ok {
		return derefStruct.dereference()
	}
	return value
}
