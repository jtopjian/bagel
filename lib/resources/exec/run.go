package exec

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jtopjian/bagel/lib/connections"
	"github.com/jtopjian/bagel/lib/utils"
)

const execRunName = "exec.Run"

type RunOpts struct {
	Command  string   `mapstructure:"cmd" required:"true"`
	Dir      string   `mapstructure:"dir"`
	Env      []string `mapstructure:"env"`
	Sudo     bool     `mapstructure:"sudo"`
	Timeout  int      `mapstructure:"timeout"`
	Unless   string   `mapstructure:"unless"`
	Internal bool

	Connection connections.Connection
	Logger     *logrus.Entry
}

func Run(input map[string]interface{}, conn connections.Connection) (*connections.RunResult, error) {
	var r RunOpts
	var result *connections.RunResult

	// validate the input
	err := utils.DecodeAndValidate(input, &r)
	if err != nil {
		return result, err
	}

	r.Connection = conn

	var internal bool
	if _, ok := input["_internal"]; ok {
		internal = true
	}

	// Configure the Logger
	var logger *logrus.Entry
	if v, ok := input["_logger"]; ok {
		if l, ok := v.(*logrus.Entry); ok {
			logger = l
		} else {
			return result, fmt.Errorf("Internal exec error: logger not set")
		}
	} else {
		logger = utils.SetLogFields(utils.GetLogger(), map[string]interface{}{
			"resource": execRunName,
		})
	}

	cmd := r.Command
	unless := r.Unless

	if r.Sudo {
		cmd = fmt.Sprintf("sudo %s", cmd)
		unless = fmt.Sprintf("sudo %s", unless)
	}

	if r.Dir != "" {
		cmd = fmt.Sprintf(`cd %s && %s`, r.Dir, cmd)
		unless = fmt.Sprintf(`cd %s && %s`, r.Dir, unless)
	}

	for _, env := range r.Env {
		cmd = fmt.Sprintf(`%s && %s`, env, cmd)
		unless = fmt.Sprintf(`%s && %s`, env, unless)
	}

	ro := connections.RunOpts{
		Command: cmd,
		Timeout: r.Timeout,
	}

	if r.Unless != "" {
		uo := connections.RunOpts{
			Command: unless,
			Timeout: r.Timeout,
		}

		if internal {
			logger.Debugf("running unless command: %s", ro.Command)
		} else {
			logger.Infof("running unless command: %s", ro.Command)
		}

		result, err = conn.RunCommand(uo)
		if result.ExitCode == 0 {
			result.Applied = false
			return result, err
		}
	}

	if internal {
		logger.Debugf("running command: %s", ro.Command)
	} else {
		logger.Infof("running command: %s", ro.Command)
	}

	result, err = conn.RunCommand(ro)
	return result, err
}

// InternalRun is like Run but takes a RunOpts argument.
// This is meant to be used internally to build more complex resources.
func InternalRun(opts RunOpts) (*connections.RunResult, error) {
	input := map[string]interface{}{
		"cmd":       opts.Command,
		"dir":       opts.Dir,
		"env":       opts.Env,
		"sudo":      opts.Sudo,
		"timeout":   opts.Timeout,
		"unless":    opts.Unless,
		"_logger":   opts.Logger,
		"_internal": true,
	}

	return Run(input, opts.Connection)
}
