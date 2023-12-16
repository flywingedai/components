package templates

const (
	InitParams = `
func initParams() {{ParamsPath}} {
	return {{ParamsPath}}{}
}
`

	BuildMocks = `
func buildMocks(t *testing.T) ({{InterfacePackage}}{{InterfaceName}}, *mocks) {
	params := initParams()

	{{MockFields}}

	return {{ComponentPackage}}New(params), convert(params)
}
`

	GetMockField = `
func mock_{{FieldName}}() {{MockPackage}}.{{MockType}}_ExpecterChain[mocks] {
	return {{MockPackage}}.ExpecterChain(func(m *mocks) *{{MockPackage}}.{{MockType}} {
		return m.{{FieldName}}
	})
}
`
)
