package file

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jtopjian/bagel/lib/connections"
	"github.com/jtopjian/bagel/lib/utils"
)

const fileDeleteName = "file.Delete"

// DeleteOpts represents options for deleting a file
type DeleteOpts struct {
	Path    string `mapstructure:"path" required:"true"`
	Timeout int    `mapstructure:"timeout"`

	Connection connections.Connection
	Logger     *logrus.Entry
}

// Delete will delete a file on a target host.
func Delete(input map[string]interface{}, conn connections.Connection) (*connections.FileResult, error) {
	var opts DeleteOpts
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
			"resource": fmt.Sprintf("%s:%s", fileDeleteName, opts.Path),
		})
	}

	if internal {
		logger.Debugf("deleting %s", opts.Path)
	} else {
		logger.Infof("deleting %s", opts.Path)
	}

	fo := connections.FileOpts{
		Path:    opts.Path,
		Timeout: opts.Timeout,
	}

	return conn.FileDelete(fo)
}

// InternalDelete is like Delete but takes a DeleteOpts argument.
// This is meant to be used internally to build more complex resources.
func InternalDelete(opts DeleteOpts) (*connections.FileResult, error) {
	input := map[string]interface{}{
		"path":      opts.Path,
		"timeout":   opts.Timeout,
		"_logger":   opts.Logger,
		"_internal": true,
	}

	return Delete(input, opts.Connection)
}
