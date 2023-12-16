package maincomponent

import (
	"bytes"
	"math/rand"
)

// 50/50 chance function
func (c *component) Chance(s rand.Source) bool {
	return c.subComponent.Chance(s)
}

// Hidden internal function
func (c *component) OtherFunction() bytes.Buffer {
	return *bytes.NewBuffer([]byte{255})
}
