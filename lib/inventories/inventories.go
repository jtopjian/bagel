package inventories

import (
	"fmt"
)

// Inventory is an interface which specifies what inventory drivers
// must implement.
type Inventory interface {
	Discover() ([]Target, error)
}

// Target represents a target returned by a driver.
type Target struct {
	Name              string
	Address           string
	ConnectionName    string
	ConnectionType    string
	ConnectionOptions map[string]interface{}
}

// New will return a target based on a given target driver.
func New(inventoryType string, options map[string]interface{}) (Inventory, error) {
	if inventoryType == "" {
		return nil, fmt.Errorf("a target type was not specified")
	}

	switch inventoryType {
	case "textfile":
		return NewTextFile(options)
	default:
		return nil, fmt.Errorf("unsupported inventory type: %s", inventoryType)
	}

	return nil, nil
}
