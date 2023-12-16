package templates

const (
	ExpecterChain = `
type {{InterfaceName}}_ExpecterChain[T any] func(*T) *{{InterfaceName}}_Expecter	
func ExpecterChain[T any](fetch func(*T) *{{InterfaceName}}) {{InterfaceName}}_ExpecterChain[T] {
	return func(m *T) *{{InterfaceName}}_Expecter {
		c := fetch(m)
		return c.EXPECT()
	}
}
`

	Chain = `
type {{InterfaceName}}_{{Method}}Chain[T any] func(*T) *{{InterfaceName}}_{{Method}}_Call

func (c {{InterfaceName}}_ExpecterChain[T]) {{Method}}({{ArgsInterface}}) {{InterfaceName}}_{{Method}}Chain[T] {
	return func(m *T) *{{InterfaceName}}_{{Method}}_Call {
		expecter := c(m)
		return expecter.{{Method}}({{ArgsShort}})
	}
}

func (c {{InterfaceName}}_{{Method}}Chain[T]) Run(run func({{Args}})) {{InterfaceName}}_{{Method}}Chain[T] {
	return func(m *T) *{{InterfaceName}}_{{Method}}_Call {
		call := c(m)
		return call.Run(run)
	}
}

func (c {{InterfaceName}}_{{Method}}Chain[T]) Return({{ReturnsArgs}}) {{InterfaceName}}_{{Method}}Chain[T] {
	return func(m *T) *{{InterfaceName}}_{{Method}}_Call {
		call := c(m)
		return call.Return({{ReturnsShort}})
	}
}

func (c {{InterfaceName}}_{{Method}}Chain[T]) RunAndReturn(run func({{Args}}){{ReturnsTypes}}) {{InterfaceName}}_{{Method}}Chain[T] {
	return func(m *T) *{{InterfaceName}}_{{Method}}_Call {
		call := c(m)
		return call.RunAndReturn(run)
	}
}
`
)
