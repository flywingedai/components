package templates

const (
	ExpecterChain = `
type {{InterfaceName}}_ExpecterChain[M any{{GenericLongAppend}}] func(*M) *{{InterfaceName}}_Expecter{{GenericShort}}
func Create_{{InterfaceName}}_ExpecterChain[M any{{GenericLongAppend}}](fetch func(*M) *{{InterfaceName}}{{GenericShort}}) {{InterfaceName}}_ExpecterChain[M{{GenericShortAppend}}] {
	return func(m *M) *{{InterfaceName}}_Expecter{{GenericShort}} {
		c := fetch(m)
		return c.EXPECT()
	}
}
`

	Chain = `
type {{InterfaceName}}_{{Method}}Chain[M any{{GenericLongAppend}}] func(*M) *{{InterfaceName}}_{{Method}}_Call{{GenericShort}}

func (c {{InterfaceName}}_ExpecterChain[M{{GenericShortAppend}}]) {{Method}}({{ArgsInterface}}) {{InterfaceName}}_{{Method}}Chain[M{{GenericShortAppend}}] {
	return func(m *M) *{{InterfaceName}}_{{Method}}_Call{{GenericShort}} {
		expecter := c(m)
		return expecter.{{Method}}({{ArgsShort}})
	}
}

func (c {{InterfaceName}}_{{Method}}Chain[M{{GenericShortAppend}}]) Run(run func({{Args}})) {{InterfaceName}}_{{Method}}Chain[M{{GenericShortAppend}}] {
	return func(m *M) *{{InterfaceName}}_{{Method}}_Call{{GenericShort}} {
		call := c(m)
		return call.Run(run)
	}
}

func (c {{InterfaceName}}_{{Method}}Chain[M{{GenericShortAppend}}]) Return({{ReturnsArgs}}) {{InterfaceName}}_{{Method}}Chain[M{{GenericShortAppend}}] {
	return func(m *M) *{{InterfaceName}}_{{Method}}_Call{{GenericShort}} {
		call := c(m)
		return call.Return({{ReturnsShort}})
	}
}

func (c {{InterfaceName}}_{{Method}}Chain[M{{GenericShortAppend}}]) RunAndReturn(run func({{Args}}){{ReturnsTypes}}) {{InterfaceName}}_{{Method}}Chain[M{{GenericShortAppend}}] {
	return func(m *M) *{{InterfaceName}}_{{Method}}_Call{{GenericShort}} {
		call := c(m)
		return call.RunAndReturn(run)
	}
}

func (c {{InterfaceName}}_ExpecterChain[M{{GenericShortAppend}}]) {{Method}}_Pointer({{ArgsInterface}}) {{InterfaceName}}_{{Method}}Chain[M{{GenericShortAppend}}] {
	return func(m *M) *{{InterfaceName}}_{{Method}}_Call{{GenericShort}} {
		expecter := c(m)
		return expecter.{{Method}}({{ArgsShortPointer}})
	}
}

func (c {{InterfaceName}}_{{Method}}Chain[M{{GenericShortAppend}}]) Return_Pointer({{ReturnsArgsPointer}}) {{InterfaceName}}_{{Method}}Chain[M{{GenericShortAppend}}] {
	return func(m *M) *{{InterfaceName}}_{{Method}}_Call{{GenericShort}} {
		call := c(m)
		return call.Return({{ReturnsShortPointer}})
	}
}
`
)
