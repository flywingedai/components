package tests

/////////////////
// DEREFERENCE //
/////////////////

type deref struct {
	dereference func() interface{}
}

// Dereference a value at runtime. Often used for pointer values
func DeRef[T any](field *T) deref {
	return deref{
		dereference: func() interface{} {
			return *field
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
