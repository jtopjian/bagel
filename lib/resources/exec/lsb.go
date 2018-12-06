package exec

import (
	"fmt"
	"regexp"

	"github.com/jtopjian/bagel/lib/resources/base"
)

type LSBInfo struct {
	DistributionID string
	Description    string
	Release        string
	Codename       string
}

func GetLSBInfo(opts base.BaseFields) (*LSBInfo, error) {
	var lsbInfo LSBInfo

	distributorRe := regexp.MustCompile("Distributor ID:\\s+(.+)\n")
	descriptionRe := regexp.MustCompile("Description:\\s+(.+)\n")
	releaseRe := regexp.MustCompile("Release:\\s+(.+)\n")
	codenameRe := regexp.MustCompile("Codename:\\s+(.+)")

	ro := RunOpts{
		Command:    "/usr/bin/lsb_release -a",
		Sudo:       opts.Sudo,
		Timeout:    opts.Timeout,
		Connection: opts.Connection,
		Logger:     opts.Logger,
	}

	result, err := InternalRun(ro)
	if err != nil {
		opts.Logger.Debug(result.Stderr)
		return nil, fmt.Errorf("unable to run lsb_info: %s", err)
	}

	if v := distributorRe.FindStringSubmatch(result.Stdout); len(v) > 1 {
		lsbInfo.DistributionID = v[1]
	}

	if v := descriptionRe.FindStringSubmatch(result.Stdout); len(v) > 1 {
		lsbInfo.Description = v[1]
	}

	if v := releaseRe.FindStringSubmatch(result.Stdout); len(v) > 1 {
		lsbInfo.Release = v[1]
	}

	if v := codenameRe.FindStringSubmatch(result.Stdout); len(v) > 1 {
		lsbInfo.Codename = v[1]
	}

	return &lsbInfo, nil
}
