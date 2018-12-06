package apt

import (
	"fmt"
	"strings"

	"github.com/jtopjian/bagel/lib/connections"
	"github.com/jtopjian/bagel/lib/resources/base"
	"github.com/jtopjian/bagel/lib/resources/exec"
	"github.com/jtopjian/bagel/lib/utils"
)

const aptPPAName = "apt.PPA"

// PPAOpts represents options for an apt.PPA resource.
type PPAOpts struct {
	base.BaseFields `mapstructure:",squash"`

	// Refresh will triger an apt-get update if set to true
	Refresh bool `default:"true"`

	fileName string
}

// PPA will perform a full state cycle for an apt.PPA.
func PPA(input map[string]interface{}, conn connections.Connection) (changed bool, err error) {
	var opts PPAOpts

	err = utils.DecodeAndValidate(input, &opts)
	if err != nil {
		return
	}

	opts.Connection = conn

	logger := utils.SetLogFields(utils.GetLogger(), map[string]interface{}{
		"resource": fmt.Sprintf("%s::%s::%s", aptPPAName, opts.Name, opts.State),
	})
	opts.Logger = logger

	exists, err := PPAExists(opts)
	if err != nil {
		return
	}

	if opts.State == "absent" {
		if exists {
			err = PPADelete(opts)
			changed = true
			return
		}

		return
	}

	if !exists {
		err = PPACreate(opts)
		changed = true
		return
	}

	return
}

// PPAExists will determine if an apt.PPA exists.
func PPAExists(opts PPAOpts) (bool, error) {
	if opts.fileName == "" {
		sourceFileName, err := aptPPASourceFileName(opts)
		if err != nil {
			return false, err
		}

		opts.fileName = "/etc/apt/sources.list.d/" + sourceFileName
		opts.Logger.Debugf("ppa file: %s", opts.fileName)
	}

	ro := exec.RunOpts{
		Command:    fmt.Sprintf(`stat "%s"`, opts.fileName),
		Sudo:       opts.Sudo,
		Timeout:    opts.Timeout,
		Connection: opts.Connection,
		Logger:     opts.Logger,
	}

	result, err := exec.InternalRun(ro)
	if err != nil {
		return false, fmt.Errorf("unable to check status of %s::%s: %s", aptPPAName, opts.Name, err)
	}

	if result.ExitCode == 0 {
		opts.Logger.Info("installed")
		return true, err
	}

	opts.Logger.Info("not installed")
	return false, nil
}

// PPACreate will create an apt.PPA resource.
func PPACreate(opts PPAOpts) error {
	ro := exec.RunOpts{
		Command:    fmt.Sprintf("apt-add-repository -y ppa:%s", opts.Name),
		Sudo:       opts.Sudo,
		Timeout:    opts.Timeout,
		Connection: opts.Connection,
		Logger:     opts.Logger,
	}

	result, err := exec.InternalRun(ro)
	if err != nil {
		opts.Logger.Debug(result.Stderr)
		return fmt.Errorf("unable to add %s::%s: %s", aptPPAName, opts.Name, err)
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("unable to add %s::%s: %s", aptPPAName, opts.Name, result.Stderr)
	}

	if opts.Refresh {
		ro.Command = "apt-get update -qq"
		result, err = exec.InternalRun(ro)
		if err != nil {
			opts.Logger.Debug(result.Stderr)
			return fmt.Errorf("unable to add %s::%s: %s", aptPPAName, opts.Name, err)
		}
	}

	return nil
}

// PPADelete will delete an apt.PPA resource.
func PPADelete(opts PPAOpts) error {
	if opts.fileName == "" {
		sourceFileName, err := aptPPASourceFileName(opts)
		if err != nil {
			return err
		}

		opts.fileName = "/etc/apt/sources.list.d/" + sourceFileName
	}

	ro := exec.RunOpts{
		Command:    fmt.Sprintf("apt-add-repository -y -r ppa:%s", opts.Name),
		Sudo:       opts.Sudo,
		Timeout:    opts.Timeout,
		Connection: opts.Connection,
		Logger:     opts.Logger,
	}

	result, err := exec.InternalRun(ro)
	if err != nil {
		opts.Logger.Debug(result.Stderr)
		return fmt.Errorf("unable to delete %s::%s: %s", aptPPAName, opts.Name, err)
	}

	if result.ExitCode != 0 {
		opts.Logger.Debug(result.Stderr)
		return fmt.Errorf("unable to delete %s::%s: %s", aptPPAName, opts.Name, err)
	}

	ro.Command = fmt.Sprintf("rm %s", opts.fileName)
	result, err = exec.InternalRun(ro)
	if err != nil {
		opts.Logger.Debug(result.Stderr)
		return fmt.Errorf("unable to delete %s::%s: %s", aptPPAName, opts.Name, err)
	}

	if result.ExitCode != 0 {
		opts.Logger.Debug(result.Stderr)
		return fmt.Errorf("unable to delete %s::%s: %s", aptPPAName, opts.Name, err)
	}

	if opts.Refresh {
		ro.Command = "apt-get update -qq"
		result, err = exec.InternalRun(ro)
		if err != nil {
			opts.Logger.Debug(result.Stderr)
			return fmt.Errorf("unable to delete %s::%s: %s", aptPPAName, opts.Name, err)
		}
	}

	opts.Logger.Info("deleted")

	return nil
}

func aptPPASourceFileName(opts PPAOpts) (string, error) {
	name := opts.Name

	lsbInfo, err := exec.GetLSBInfo(opts.BaseFields)
	if err != nil {
		return "", nil
	}

	distro := fmt.Sprintf("-%s-", strings.ToLower(lsbInfo.DistributionID))
	release := strings.ToLower(lsbInfo.Codename)

	name = strings.Replace(name, "/", distro, -1)
	name = strings.Replace(name, ":", "-", -1)
	name = strings.Replace(name, ".", "_", -1)

	name = fmt.Sprintf("%s-%s.list", name, release)

	return name, nil
}
