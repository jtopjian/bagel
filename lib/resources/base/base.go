package base

import (
	"github.com/sirupsen/logrus"

	"github.com/jtopjian/bagel/lib/connections"
)

// BaseFields represents fields which are
// standard to all resources.
type BaseFields struct {
	// Name is the name of the resource. The value
	// will differ from resource to resource.
	Name string `mapstructure:"name" required:"true"`

	// State represents the state of the resource.
	// It can either be "present", "absent", "latest",
	// or a version number
	State string `mapstructure:"state" default:"present"`

	// Sudo is if the command requires sudo to run.
	Sudo bool `mapstructure:"sudo"`

	// Timeout is a timeout for the command.
	Timeout int `mapstructure:"timeout"`

	// Connection represents an internal connection to use
	// to execute commands on the host.
	Connection connections.Connection

	// Logger represents an internal logger.
	Logger *logrus.Entry
}
