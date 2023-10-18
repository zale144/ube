package product

import (
	"testing"

	"github.com/zale144/ube/libs/testengine"
)

func TestProduct(t *testing.T) {
	testengine.EventHandler(t, &Product{})
}
