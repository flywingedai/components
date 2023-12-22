package templates

const Interface = `
type %s%s interface {
%s
}
`

const Method = "\t%s(%s) %s"

const New = `
func New{{GenericLong}}(p Params{{GenericShort}}) {{Interface}}{{GenericShort}} {
	return p.Convert()
}
`
