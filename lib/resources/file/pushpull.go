package file

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/jtopjian/bagel/lib/connections"
	"github.com/jtopjian/bagel/lib/resources/exec"
	"github.com/jtopjian/bagel/lib/utils"
)

// PushPullOpts represents options for a push/pull action.
type PushPullOpts struct {
	Source      string `mapstructure:"source" required:"true"`
	Destination string `mapstructure:"destination" required:"true"`
	UID         int    `mapstructure:"uid"`
	GID         int    `mapstructure:"gid"`
	Mode        int    `mapstructure:"mode"`
	Timeout     int    `mapstructure:"timeout"`

	Connection connections.Connection
	Logger     *logrus.Entry
}

func pushPull(input map[string]interface{}, conn connections.Connection, action string) (*connections.FileResult, error) {
	var opts PushPullOpts
	var result *connections.FileResult

	// validate the input
	err := utils.DecodeAndValidate(input, &opts)
	if err != nil {
		return result, err
	}

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
			"resource": fmt.Sprintf("file.%s:%s", strings.Title(action), opts.Source),
		})
	}

	if internal {
		logger.Debugf("%s %s => %s", action, opts.Source, opts.Destination)
	} else {
		logger.Infof("%s %s => %s", action, opts.Source, opts.Destination)
	}

	cfo := connections.CopyFileOpts{
		Source:      opts.Source,
		Destination: opts.Destination,
		UID:         opts.UID,
		GID:         opts.GID,
		Mode:        os.FileMode(opts.Mode),
		Timeout:     opts.Timeout,
	}

	switch action {
	case "push":
		return conn.FileUpload(cfo)
	case "pull":
		return conn.FileDownload(cfo)
	}

	return nil, nil
}

// Push will push/upload a file to a target host.
func Push(input map[string]interface{}, conn connections.Connection) (*connections.FileResult, error) {
	return pushPull(input, conn, "push")
}

// InternalPush is like Push, but takes a PushPullOpts argument.
// This is meant to be used internally to build more complex resources.
func InternalPush(opts PushPullOpts) (*connections.FileResult, error) {
	input := map[string]interface{}{
		"source":      opts.Source,
		"destination": opts.Destination,
		"uid":         opts.UID,
		"gid":         opts.GID,
		"mode":        opts.Mode,
		"timeout":     opts.Timeout,
		"_logger":     opts.Logger,
		"_internal":   true,
	}

	return Push(input, opts.Connection)
}

// PushAndMove is a convenience function to upload/push a file
// to a target host and then move it to a secondary location. This
// is useful for cases when `sudo` is required.
func PushAndMove(pushPullOpts PushPullOpts, runOpts exec.RunOpts, finalDestination string) (*connections.FileResult, error) {
	input := map[string]interface{}{
		"source":      pushPullOpts.Source,
		"destination": pushPullOpts.Destination,
		"uid":         pushPullOpts.UID,
		"gid":         pushPullOpts.GID,
		"mode":        pushPullOpts.Mode,
		"timeout":     pushPullOpts.Timeout,
		"_logger":     pushPullOpts.Logger,
		"_internal":   true,
	}

	result, err := Push(input, pushPullOpts.Connection)
	if err != nil {
		return nil, err
	}

	if !result.Success {
		return nil, fmt.Errorf("unable to upload file to %s", pushPullOpts.Destination)
	}

	runOpts.Command = fmt.Sprintf(`mv "%s" "%s"`, pushPullOpts.Destination, finalDestination)
	rr, err := exec.InternalRun(runOpts)
	if err != nil {
		return nil, err
	}

	if rr.ExitCode != 0 {
		return nil, fmt.Errorf("unable to upload file to %s: %s", pushPullOpts.Destination, rr.Stderr)
	}

	return result, nil
}

// Pull will pull/download a file from a target host.
func Pull(input map[string]interface{}, conn connections.Connection) (*connections.FileResult, error) {
	return pushPull(input, conn, "pull")
}

// InternalPull is like Pull, but takes a PushPullOpts argument.
// This is meant to be used internally to build more complex resources.
func InternalPull(opts PushPullOpts) (*connections.FileResult, error) {
	input := map[string]interface{}{
		"source":      opts.Source,
		"destination": opts.Destination,
		"uid":         opts.UID,
		"gid":         opts.GID,
		"mode":        opts.Mode,
		"timeout":     opts.Timeout,
		"_logger":     opts.Logger,
		"_internal":   true,
	}

	return Pull(input, opts.Connection)
}
