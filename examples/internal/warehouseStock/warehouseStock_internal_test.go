package warehousestock

import (
	"testing"

	"github.com/zale144/ube/libs/testengine"
)

func TestWarehouseStock(t *testing.T) {
	testengine.EventHandler(t, &UBEModel{})
}
