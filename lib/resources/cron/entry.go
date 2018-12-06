package cron

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

	"github.com/mitchellh/mapstructure"
)

const cronEntryName = "cron.Entry"

// EntryOpts represents options for a cron.Entry action.
type EntryOpts struct {
	base.BaseFields `mapstructure:",squash"`

	// User is the user who owns the cron entry.
	User string `default:"root"`

	// Command is the command which cron will run.
	Command string `required:"true"`

	// Minute is the minute field of the cron entry.
	Minute string `default:"*"`

	// Hour is the hour field of the cron entry.
	Hour string `default:"*"`

	// DayOfMonth is the day of the month field of the cron entry.
	DayOfMonth string `default:"*"`

	// Month is the month field of the cron entry.
	Month string `default:"*"`

	// DayOfWeek is the day of the week field of the cron entry.
	DayOfWeek string `default:"*"`
}

// entry returns the formatted cron entry.
func (opts EntryOpts) entry() string {
	entry := fmt.Sprintf(`%s %s %s %s %s %s # %s`,
		opts.Minute, opts.Hour, opts.DayOfMonth, opts.Month,
		opts.DayOfWeek, opts.Command, opts.Name)

	return entry
}

// Entry will perform a full state cycle for a cron.Entry.
func Entry(input map[string]interface{}, conn connections.Connection) (changed bool, err error) {
	var opts EntryOpts

	config := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &opts,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return
	}

	err = decoder.Decode(input)
	if err != nil {
		return
	}

	err = utils.ValidateTags(&opts)
	if err != nil {
		return
	}

	opts.Connection = conn

	logger := utils.SetLogFields(utils.GetLogger(), map[string]interface{}{
		"resource": fmt.Sprintf("%s::%s::%s", cronEntryName, opts.Name, opts.State),
	})
	opts.Logger = logger

	exists, err := EntryExists(opts)
	if err != nil {
		return
	}

	if opts.State == "absent" {
		if exists {
			err = EntryDelete(opts)
			changed = true
			return
		}

		return
	}

	if !exists {
		err = EntryCreate(opts)
		changed = true
		return
	}

	return
}

// EntryExists will determine if a cron.Entry exists.
func EntryExists(opts EntryOpts) (bool, error) {
	entries, err, stderr := getCronEntries(opts)
	if err != nil {
		return false, fmt.Errorf("unable to check status of %s::%s: %s", cronEntryName, opts.Name, err)
	}

	if stderr != nil {
		opts.Logger.Info("not installed")
		return false, nil
	}

	var exists bool
	for _, line := range entries {
		if line == opts.entry() {
			exists = true
		}
	}

	if exists {
		opts.Logger.Info("installed")
		return true, nil
	}

	opts.Logger.Info("not installed")
	return false, nil
}

// EntryCreate will create a cron.entry.
func EntryCreate(opts EntryOpts) error {
	entries, err, _ := getCronEntries(opts)
	if err != nil {
		return fmt.Errorf("unable to add %s::%s: %s", cronEntryName, opts.Name, err)
	}

	var newEntries []string
	var added bool
	for _, line := range entries {
		if strings.Contains(line, fmt.Sprintf(`# %s`, opts.Name)) {
			line = opts.entry()
			added = true
		}
		newEntries = append(newEntries, line)
	}

	if !added {
		newEntries = append(newEntries, opts.entry())
	}

	newEntries = append(newEntries, "\n")

	if err := pushCronEntries(opts, newEntries); err != nil {
		return fmt.Errorf("unable to add %s::%s: %s", cronEntryName, opts.Name, err)
	}

	return nil
}

// EntryDelete will delete a cron.Entry.
func EntryDelete(opts EntryOpts) error {
	entries, err, stderr := getCronEntries(opts)
	if err != nil {
		return fmt.Errorf("unable to delete %s::%s: %s", cronEntryName, opts.Name, err)
	}

	if stderr != nil {
		return fmt.Errorf("unable to add %s::%s: %s", cronEntryName, opts.Name, stderr)
	}

	var newEntries []string
	for _, line := range entries {
		if line != opts.entry() {
			newEntries = append(newEntries, line)
		}
	}

	newEntries = append(newEntries, "\n")

	if err := pushCronEntries(opts, newEntries); err != nil {
		return fmt.Errorf("unable to delete %s::%s: %s", cronEntryName, opts.Name, err)
	}

	return nil
}

// getCronEntries returns the cron entries from a remote host.
func getCronEntries(opts EntryOpts) ([]string, error, error) {
	ro := exec.RunOpts{
		Command:    fmt.Sprintf("crontab -u %s -l", opts.User),
		Sudo:       opts.Sudo,
		Timeout:    opts.Timeout,
		Connection: opts.Connection,
		Logger:     opts.Logger,
	}

	result, err := exec.InternalRun(ro)
	if err != nil {
		opts.Logger.Debug(result.Stderr)
		return nil, nil, err
	}

	if result.ExitCode != 0 {
		return nil, nil, fmt.Errorf("%s", result.Stderr)
	}

	return strings.Split(result.Stdout, "\n"), nil, nil
}

// pushCronEntries pushes new entries to a remote host.
func pushCronEntries(opts EntryOpts, entries []string) error {
	tmpfile, err := ioutil.TempFile("/tmp", "cron.entry")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(strings.Join(entries, "\n"))); err != nil {
		return err
	}

	if err := tmpfile.Close(); err != nil {
		return err
	}

	ppo := file.PushPullOpts{
		Source:      tmpfile.Name(),
		Destination: tmpfile.Name(),
		Connection:  opts.Connection,
		Logger:      opts.Logger,
	}

	if _, err := file.InternalPush(ppo); err != nil {
		return err
	}

	ro := exec.RunOpts{
		Command:    fmt.Sprintf(`crontab -u %s %s`, opts.User, tmpfile.Name()),
		Sudo:       opts.Sudo,
		Timeout:    opts.Timeout,
		Connection: opts.Connection,
		Logger:     opts.Logger,
	}

	result, err := exec.InternalRun(ro)
	if err != nil {
		opts.Logger.Debug(result.Stderr)
		return err
	}

	if result.ExitCode != 0 {
		return fmt.Errorf("%s", result.Stderr)
	}

	dfo := file.DeleteOpts{
		Path:       tmpfile.Name(),
		Connection: opts.Connection,
		Logger:     opts.Logger,
	}

	_, err = file.InternalDelete(dfo)
	if err != nil {
		return err
	}

	return nil
}
