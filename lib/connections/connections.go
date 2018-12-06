package connections

import (
	"fmt"
)

// New will return a connection based on a given connection type.
func New(connType string, options map[string]interface{}) (Connection, error) {
	if connType == "" {
		return nil, fmt.Errorf("a connection type was not specified")
	}

	switch connType {
	case "local":
		return NewLocal(options)
	case "ssh":
		return NewSSH(options)
	default:
		return nil, fmt.Errorf("unsupported connection type: %s", connType)
	}

	return nil, nil
}
