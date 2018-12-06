package apt

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jtopjian/bagel/lib/connections"
	"github.com/jtopjian/bagel/lib/resources/base"
	"github.com/jtopjian/bagel/lib/resources/exec"
	"github.com/jtopjian/bagel/lib/resources/file"
	"github.com/jtopjian/bagel/lib/utils"
)

const aptSourceName = "apt.Source"

// SourceOpts represents options for an apt.Source resource.
type SourceOpts struct {
	base.BaseFields `mapstructure:",squash"`

	URI          string `required:"true"`
	Distribution string `required:"true"`
	Component    string
	IncludeSrc   bool
	Refresh      bool `default:"true"`
}

// Source will perform a full state cycle for an apt source entry.
func Source(input map[string]interface{}, conn connections.Connection) (changed bool, err error) {
	var opts SourceOpts

	err = utils.DecodeAndValidate(input, &opts)
	if err != nil {
		return
	}

	opts.Connection = conn

	logger := utils.SetLogFields(utils.GetLogger(), map[string]interface{}{
		"resource": fmt.Sprintf("%s::%s::%s", aptSourceName, opts.Name, opts.State),
	})
	opts.Logger = logger

	exists, err := SourceExists(opts)
	if err != nil {
		return
	}

	if opts.State == "absent" {
		if exists {
			err = SourceDelete(opts)
			changed = true
			return
		}

		return
	}

	if !exists {
		err = SourceCreate(opts)
		changed = true
		return
	}

	return
}

// SourceExists will determine if an apt.Source exists.
func SourceExists(opts SourceOpts) (bool, error) {
	path := fmt.Sprintf("/etc/apt/sources.list.d/%s.list", opts.Name)
	entry := fmt.Sprintf("deb %s %s %s", opts.URI, opts.Distribution, opts.Component)
	srcEntry := fmt.Sprintf("deb-src %s %s %s", opts.URI, opts.Distribution, opts.Component)

	ro := exec.RunOpts{
		Command:    fmt.Sprintf(`cat "%s"`, path),
		Sudo:       opts.Sudo,
		Timeout:    opts.Timeout,
		Connection: opts.Connection,
		Logger:     opts.Logger,
	}

	result, err := exec.InternalRun(ro)
	if err != nil {
		return false, fmt.Errorf("unable to check status of %s::%s: %s", aptSourceName, opts.Name, err)
	}

	if result.ExitCode != 0 {
		opts.Logger.Info("not installed")
		return false, nil
	}

	var exists bool
	var srcExists bool
	for _, line := range strings.Split(result.Stdout, "\n") {
		if line == entry {
			exists = true
		}

		if line == srcEntry {
			srcExists = true
		}
	}

	if exists {
		if opts.IncludeSrc && !srcExists {
			opts.Logger.Info("not installed")
			return false, nil
		}

		opts.Logger.Info("exists")
		return true, nil
	}

	opts.Logger.Info("not installed")
	return false, nil
}

// SourceCreate will create an apt.Source file.
func SourceCreate(opts SourceOpts) error {
	path := fmt.Sprintf("/etc/apt/sources.list.d/%s.list", opts.Name)
	entry := fmt.Sprintf("deb %s %s %s", opts.URI, opts.Distribution, opts.Component)
	srcEntry := fmt.Sprintf("\ndeb-src %s %s %s", opts.URI, opts.Distribution, opts.Component)

	tmpfile, err := ioutil.TempFile("/tmp", "apt.source")
	if err != nil {
		return fmt.Errorf("unable to add %s::%s: %s", aptSourceName, opts.Name, err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(entry)); err != nil {
		return fmt.Errorf("unable to add %s::%s: %s", aptSourceName, opts.Name, err)
	}

	if opts.IncludeSrc {
		if _, err := tmpfile.Write([]byte(srcEntry)); err != nil {
			return fmt.Errorf("unable to add %s::%s: %s", aptSourceName, opts.Name, err)
		}
	}

	if err := tmpfile.Close(); err != nil {
		return fmt.Errorf("unable to add %s::%s: %s", aptSourceName, opts.Name, err)
	}

	ppo := file.PushPullOpts{
		Source:      tmpfile.Name(),
		Destination: tmpfile.Name(),
		Connection:  opts.Connection,
		Logger:      opts.Logger,
	}

	ro := exec.RunOpts{
		Sudo:       opts.Sudo,
		Timeout:    opts.Timeout,
		Connection: opts.Connection,
		Logger:     opts.Logger,
	}

	if _, err := file.PushAndMove(ppo, ro, path); err != nil {
		return fmt.Errorf("unable to add %s::%s: %s", aptSourceName, opts.Name, err)
	}

	if opts.Refresh {
		ro.Command = "apt-get update -qq"
		result, err := exec.InternalRun(ro)
		if err != nil {
			opts.Logger.Debug(result.Stderr)
			return fmt.Errorf("unable to add %s::%s: %s", aptSourceName, opts.Name, err)
		}
	}

	return nil
}

// Delete will delete an apt.source file.
func SourceDelete(opts SourceOpts) error {
	path := fmt.Sprintf("/etc/apt/sources.list.d/%s.list", opts.Name)
	ro := exec.RunOpts{
		Command:    fmt.Sprintf(`rm "%s"`, path),
		Sudo:       opts.Sudo,
		Timeout:    opts.Timeout,
		Connection: opts.Connection,
		Logger:     opts.Logger,
	}

	result, err := exec.InternalRun(ro)
	if err != nil {
		opts.Logger.Debug(result.Stderr)
		return fmt.Errorf("unable to delete %s::%s: %s", aptSourceName, opts.Name, err)
	}

	if result.ExitCode != 0 {
		opts.Logger.Debug(result.Stderr)
		return fmt.Errorf("unable to delete %s::%s: %s", aptSourceName, opts.Name, err)
	}

	if opts.Refresh {
		ro.Command = "apt-get update -qq"
		result, err = exec.InternalRun(ro)
		if err != nil {
			opts.Logger.Debug(result.Stderr)
			return fmt.Errorf("unable to delete %s::%s: %s", aptSourceName, opts.Name, err)
		}
	}

	return nil
}
