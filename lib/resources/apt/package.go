package apt

import (
	"fmt"
	"regexp"

	"github.com/jtopjian/bagel/lib/connections"
	"github.com/jtopjian/bagel/lib/resources/base"
	"github.com/jtopjian/bagel/lib/resources/exec"
	"github.com/jtopjian/bagel/lib/utils"
)

const aptPkgName = "apt.Package"

// PackageOpts represents options for an apt.Package resource.
type PackageOpts struct {
	base.BaseFields `mapstructure:",squash"`
}

// Package will perform a full state cycle for an apt package.
func Package(input map[string]interface{}, conn connections.Connection) (changed bool, err error) {
	var opts PackageOpts

	err = utils.DecodeAndValidate(input, &opts)
	if err != nil {
		return
	}

	opts.Connection = conn

	logger := utils.SetLogFields(utils.GetLogger(), map[string]interface{}{
		"resource": fmt.Sprintf("%s::%s::%s", aptPkgName, opts.Name, opts.State),
	})
	opts.Logger = logger

	exists, err := PackageExists(opts)
	if err != nil {
		return
	}

	if opts.State == "absent" {
		if exists {
			err = PackageDelete(opts)
			changed = true
			return
		}

		return
	}

	if !exists || opts.State == "latest" {
		err = PackageCreate(opts)
		changed = true
		return
	}

	return
}

// PackageExists will determine if an apt package exists.
func PackageExists(opts PackageOpts) (bool, error) {
	e := exec.RunOpts{
		Command:    fmt.Sprintf("apt-cache policy %s", opts.Name),
		Sudo:       opts.Sudo,
		Timeout:    opts.Timeout,
		Connection: opts.Connection,
		Logger:     opts.Logger,
	}

	result, err := exec.InternalRun(e)
	if err != nil {
		opts.Logger.Debug(result.Stderr)
		return false, fmt.Errorf("unable to check status of %s::%s: %s", aptPkgName, opts.Name, err)
	}

	if result.Stdout == "" {
		return false, fmt.Errorf("no such package")
	}

	installedVersion, _ := aptPkgParseAptCache(result.Stdout)

	switch opts.State {
	case "present", "absent", "":
		switch installedVersion {
		case "(none)":
			opts.Logger.Info("not installed")
			return false, nil
		default:
			opts.Logger.Info("installed")
			return true, nil
		}
	}

	if opts.State != installedVersion {
		opts.Logger.Info("will be installed")
		return false, nil
	}

	opts.Logger.Info("installed")
	return true, nil
}

func PackageCreate(opts PackageOpts) error {
	e := exec.RunOpts{
		Sudo:       opts.Sudo,
		Timeout:    opts.Timeout,
		Connection: opts.Connection,
		Logger:     opts.Logger,
	}

	e.Env = []string{
		"DEBIAN_FRONTEND=noninteractive",
		"APT_LISTBUGS_FRONTEND=none",
		"APT_LISTCHANGES_FRONTEND=none",
	}

	var createArgs string
	if opts.State != "present" && opts.State != "latest" {
		createArgs = fmt.Sprintf("%s=%s", opts.Name, opts.State)
	} else {
		createArgs = opts.Name
	}

	e.Command = fmt.Sprintf(
		"apt-get install -y --allow-downgrades --allow-remove-essential "+
			"--allow-change-held-packages -o DPkg::Options::=--force-confold %s",
		createArgs)

	opts.Logger.Info("installing")

	result, err := exec.InternalRun(e)
	if err != nil {
		opts.Logger.Debug(result.Stderr)
		return fmt.Errorf("unable to install %s %s: %s", aptPkgName, opts.Name, err)
	}

	opts.Logger.Info("installed")
	return nil
}

func PackageDelete(opts PackageOpts) error {
	e := exec.RunOpts{
		Sudo:       opts.Sudo,
		Timeout:    opts.Timeout,
		Connection: opts.Connection,
		Logger:     opts.Logger,
	}

	e.Env = []string{
		"DEBIAN_FRONTEND=noninteractive",
		"APT_LISTBUGS_FRONTEND=none",
		"APT_LISTCHANGES_FRONTEND=none",
	}

	e.Command = fmt.Sprintf("apt-get purge -q -y %s", opts.Name)

	opts.Logger.Info("removing")

	result, err := exec.InternalRun(e)
	if err != nil {
		opts.Logger.Debug(result.Stderr)
		return fmt.Errorf("unable to remove %s %s: %s", aptPkgName, opts.Name, err)
	}
	opts.Logger.Debug(result.Stderr)

	opts.Logger.Info("removed")
	return nil
}

// aptPkgParseAptCache is an internal function that will parse the
// output of apt-cache policy and return the version information.
func aptPkgParseAptCache(stdout string) (installed, candidate string) {
	installedRe := regexp.MustCompile("Installed: (.+)\n")
	candidateRe := regexp.MustCompile("Candidate: (.+)\n")

	if v := installedRe.FindStringSubmatch(stdout); len(v) > 1 {
		installed = v[1]
	}

	if v := candidateRe.FindStringSubmatch(stdout); len(v) > 1 {
		candidate = v[1]
	}

	return
}
