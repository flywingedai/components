package maincomponent

import (
	"bytes"
	"math/rand"

	"github.com/flywingedai/components/example/subcomponent"
)

type Params struct {
	InternalField int
	SubComponent  subcomponent.SubComponent
}

type component struct {
	/*
		generate::components
		interfaceName::MainComponent
		config::.mockery.yml
		blackbox::true
	*/
	computedField int
	internalField int
	subComponent  subcomponent.SubComponent `pkg:"-"`
}

func (p *Params) Convert() *component {
	computedField := 2 * p.InternalField
	return &component{
		computedField: computedField,
		internalField: p.InternalField,
		subComponent:  p.SubComponent,
	}
}

// Code below was generated by components. DO NOT EDIT.
// Component version: v0.1.0

type MainComponent interface {
	Chance(s rand.Source) bool
	OtherFunction() bytes.Buffer
}

func New(p Params) MainComponent {
	return p.Convert()
}