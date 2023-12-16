package templates

const Interface = `
type %s interface {
%s
}
`

const Method = "\t%s(%s) (%s)"

const New = `
func New(p Params) %s {
	return p.Convert()
}
`
