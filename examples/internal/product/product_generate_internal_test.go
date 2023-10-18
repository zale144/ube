//go:build generatetests

package product

import (
	"testing"

	"github.com/zale144/ube/libs/testengine"
)

/*
If you want to run the test, just rename the go:build above into no:build, and change it back afterwards.
*/
func TestProductGenerate(t *testing.T) {
	tCase := testengine.Case{
		Entity:            &Product{},
		EventNames:        []string{create, update},
		NumMessages:       1,
		RecordsPerMessage: 1,
		WithNegative:      true,
	}

	testengine.EventHandlerInit(t, tCase)
}
