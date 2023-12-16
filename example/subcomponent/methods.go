package subcomponent

import (
	"math/rand"
)

// 50/50 chance function
func (c *component) Chance(s rand.Source) bool {
	generator := c.getRand(s)
	return generator.Intn(c.computedField) > c.internalField
}

// Hidden internal function
func (c *component) getRand(s rand.Source) *rand.Rand {
	return rand.New(s)
}
