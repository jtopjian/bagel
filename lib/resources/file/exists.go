package file

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jtopjian/bagel/lib/connections"
	"github.com/jtopjian/bagel/lib/utils"
)

const fileExistsName = "file.Exists"

// ExistsOpts represents options for checking if a file exists.
type ExistsOpts struct {
	Path    string `mapstructure:"path" required:"true"`
	Timeout int    `mapstructure:"timeout"`

	Connection connections.Connection
	Logger     *logrus.Entry
}

// Exists will determine if a file exists on a target host.
func Exists(input map[string]interface{}, conn connections.Connection) (*connections.FileResult, error) {
	var opts ExistsOpts
	var result *connections.FileResult

	// validate the input
	err := utils.DecodeAndValidate(input, &opts)
	if err != nil {
		return result, err
	}

	opts.Connection = conn

	var internal bool
	if _, ok := input["_internal"]; ok {
		internal = true
	}

	var logger *logrus.Entry
	if v, ok := input["_logger"]; ok {
		if l, ok := v.(*logrus.Entry); ok {
			logger = l
		} else {
			return result, fmt.Errorf("Internal file error: logger not set")
		}
	} else {
		logger = utils.SetLogFields(utils.GetLogger(), map[string]interface{}{
			"resource": fmt.Sprintf("%s:%s", fileExistsName, opts.Path),
		})
	}

	if internal {
		logger.Debugf("checking existence of %s", opts.Path)
	} else {
		logger.Infof("checking existence of %s", opts.Path)
	}

	fo := connections.FileOpts{
		Path:    opts.Path,
		Timeout: opts.Timeout,
	}

	return conn.FileInfo(fo)
}

// InternalExists is like Exists, but takes an ExistsOpts argument.
// This is meant to be used internally to build more exomplex resources.
func InternalExists(opts ExistsOpts, conn connections.Connection) (*connections.FileResult, error) {
	input := map[string]interface{}{
		"path":      opts.Path,
		"timeout":   opts.Timeout,
		"_logger":   opts.Logger,
		"_internal": true,
	}

	return Exists(input, opts.Connection)
}
