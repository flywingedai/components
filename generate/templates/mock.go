package templates

const (
	InitParams = `
func initParams{{GenericLong}}() {{ParamsPath}}{{GenericShort}} {
	return {{ParamsPath}}{{GenericShort}}{}
}
`

	BuildMocks = `
func buildMocks{{GenericLong}}(t *testing.T) ({{InterfacePackage}}{{InterfaceName}}{{GenericShort}}, *mocks{{GenericShort}}) {
	params := initParams{{GenericShort}}()

	{{MockFields}}

	return {{ComponentPackage}}New(params), convert(params)
}
`

	GetMockField = `
func mock_{{FieldName}}() {{MockPackage}}.{{MockType}}_ExpecterChain[mocks{{GenericShortAppend}}] {
	return {{MockPackage}}.Create_{{MockType}}_ExpecterChain(func(m *mocks) *{{MockPackage}}.{{MockType}}{{GenericShort}} {
		return m.{{FieldName}}
	})
}
`
)
