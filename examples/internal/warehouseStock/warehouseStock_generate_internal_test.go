//go:build generatetests

package warehousestock

import (
	"testing"

	"github.com/zale144/ube/libs/testengine"
)

/*
If you want to run the test, just rename the go:build above into no:build, and change it back afterwards.
*/
func TestWarehouseStockGenerate(t *testing.T) {
	tc := testengine.Case{
		Entity:            &UBEModel{},
		Feed:              &Feed{},
		NumMessages:       3,
		RecordsPerMessage: 2,
		WithNegative:      true,
	}
	testengine.EventHandlerInit(t, tc)
}
