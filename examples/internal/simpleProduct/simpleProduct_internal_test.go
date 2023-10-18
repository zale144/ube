package product

import (
	"testing"

	"github.com/zale144/ube/libs/testengine"
)

func TestSimpleProduct(t *testing.T) {
	testengine.EventHandler(t, &SimpleProduct{})
}
