package store

import (
	"fmt"

	"github.com/zale144/ube/model"
)

// Store is the example Store structure
type Store struct {
	EntityKey
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

type EntityKey struct {
	ID int `json:"id,omitempty"`
}

func (e EntityKey) PK() string {
	return fmt.Sprint(e.ID)
}

// GetKey returns the custom composed key to the entity
func (s Store) GetKey() model.Key {
	return s.EntityKey
}
