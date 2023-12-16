package maincomponent_test

import (
	"math/rand"
	"testing"

	"github.com/flywingedai/components/tests"
	"github.com/stretchr/testify/mock"
)

func TestTest(t *testing.T) {
	tester := tests.NewTesterWithoutData(buildMocks)
	test := tester.Options.
		Mock(mock_subComponent().Chance(mock.Anything).Return(true)).
		SetInputByValue(0, rand.NewSource(rand.Int63())).
		CreateMethodTest("Chance", "blargh")

	tester.AddTests(test)
	tester.Test(t)
}
