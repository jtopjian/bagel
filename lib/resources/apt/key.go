package apt

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/crypto/openpgp"

	"github.com/jtopjian/bagel/lib/connections"
	"github.com/jtopjian/bagel/lib/resources/base"
	"github.com/jtopjian/bagel/lib/resources/exec"
	"github.com/jtopjian/bagel/lib/resources/file"
	"github.com/jtopjian/bagel/lib/utils"
)

const aptKeyName = "apt.Key"

// KeyOpts represents options for an apt.Key resource.
type KeyOpts struct {
	base.BaseFields `mapstructure:",squash"`

	// KeyServer is an optional remote server to obtain the key from.
	// If KeyServer is not used, RemoteKeyFile must be used.
	KeyServer string

	// RemoteKeyFile is the URL to a public key.
	// If RemoteKeyFile is not used, KeyServer must be used.
	RemoteKeyFile string
}

// Key will perform a full state cycle for an apt key.
func Key(input map[string]interface{}, conn connections.Connection) (changed bool, err error) {
	var opts KeyOpts

	err = utils.DecodeAndValidate(input, &opts)
	if err != nil {
		return
	}

	if opts.KeyServer == "" && opts.RemoteKeyFile == "" {
		err = fmt.Errorf("%s: one of key_server or remote_key_file must be specified", opts.Name)
		return
	}

	opts.Connection = conn

	logger := utils.SetLogFields(utils.GetLogger(), map[string]interface{}{
		"resource": fmt.Sprintf("%s::%s::%s", aptKeyName, opts.Name, opts.State),
	})
	opts.Logger = logger

	exists, err := KeyExists(opts)
	if err != nil {
		return
	}

	if opts.State == "absent" {
		if exists {
			err = KeyDelete(opts)
			changed = true
			return
		}

		return
	}

	if !exists {
		err = KeyCreate(opts)
		changed = true
		return
	}

	return
}

// KeyExists will determine if a key exists
func KeyExists(opts KeyOpts) (bool, error) {
	ro := exec.RunOpts{
		Command:    fmt.Sprintf("apt-key export %s", opts.Name),
		Sudo:       opts.Sudo,
		Timeout:    opts.Timeout,
		Connection: opts.Connection,
		Logger:     opts.Logger,
	}

	result, err := exec.InternalRun(ro)
	if err != nil {
		opts.Logger.Debug(result.Stderr)
		return false, fmt.Errorf("unable to check status of %s::%s: %s", aptKeyName, opts.Name, err)
	}

	if result.Stdout == "" {
		opts.Logger.Info("not installed")
		return false, nil
	}

	/*
		_, err = aptKeyGetName(rr.Stdout)
		if err != nil {
			return false, err
		}
	*/

	opts.Logger.Info("installed")
	return true, nil
}

// KeyCreate will create a key via apt-key.
func KeyCreate(opts KeyOpts) error {
	ro := exec.RunOpts{
		Sudo:       opts.Sudo,
		Timeout:    opts.Timeout,
		Connection: opts.Connection,
		Logger:     opts.Logger,
	}

	if opts.RemoteKeyFile != "" {
		k, err := aptKeyGetRemoteKeyFile(opts.RemoteKeyFile)
		if err != nil {
			return err
		}

		tmpfile, err := ioutil.TempFile("/tmp", "apt.key")
		if err != nil {
			return err
		}
		defer os.Remove(tmpfile.Name())

		if _, err = tmpfile.Write([]byte(k)); err != nil {
			return err
		}

		if err = tmpfile.Close(); err != nil {
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

		ro.Command = fmt.Sprintf("apt-key add %s", tmpfile.Name())
		result, err := exec.InternalRun(ro)
		if err != nil {
			return err
		}

		if result.ExitCode != 0 {
			opts.Logger.Debug(result.Stderr)
			return fmt.Errorf("unable to add key: %s", err)
		}

		fdo := file.DeleteOpts{
			Path:       tmpfile.Name(),
			Connection: opts.Connection,
			Logger:     opts.Logger,
		}
		fr, err := file.InternalDelete(fdo)
		if err != nil {
			return err
		}

		if !fr.Success {
			return fmt.Errorf("unable to delete temporary key from remote host: %s", tmpfile.Name())
		}
	}

	if opts.KeyServer != "" {
		ro.Command = fmt.Sprintf("apt-key adv --keyserver %s --recv-keys %s",
			opts.KeyServer, opts.Name)

		result, err := exec.InternalRun(ro)
		if err != nil {
			return err
		}

		if result.Stderr != "" {
			return fmt.Errorf("unable to add key: %s", err)
		}
	}

	opts.Logger.Info("installed")
	return nil
}

// Delete deletes a key managed by apt.key.
func KeyDelete(opts KeyOpts) error {
	ro := exec.RunOpts{
		Command:    fmt.Sprintf("apt-key del %s", opts.Name),
		Sudo:       opts.Sudo,
		Timeout:    opts.Timeout,
		Connection: opts.Connection,
		Logger:     opts.Logger,
	}

	result, err := exec.InternalRun(ro)
	if err != nil {
		return err
	}

	if result.ExitCode != 0 {
		opts.Logger.Debug(result.Stderr)
		return fmt.Errorf("unable to delete key: %s", err)
	}

	opts.Logger.Info("deleted")
	return nil
}

// aptKeyGetRemoteKeyFile is an internal function that will
// download a key located at a remote URL.
func aptKeyGetRemoteKeyFile(v string) (key string, err error) {
	res, err := http.Get(v)
	if err != nil {
		return
	}

	k, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	key = string(k)

	return
}

// aptKeyGetShortID is an internal function that will print the
// short key ID of a public key.
func aptKeyGetShortID(key string) (fingerprint string, err error) {
	el, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(key))
	if err != nil {
		return
	}

	if len(el) == 0 {
		err = fmt.Errorf("Error determining fingerprint of key")
		return
	}

	fingerprint = el[0].PrimaryKey.KeyIdShortString()

	return
}

// aptKeyGetName is an internal function that will get the
// maintainer name of a public key.
func aptKeyGetName(key string) (name string, err error) {
	el, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(key))
	if err != nil {
		return
	}

	if len(el) == 0 {
		err = fmt.Errorf("Error determining userid of key")
		return
	}

	identities := el[0].Identities
	for k, _ := range identities {
		if name == "" {
			name = k
		}
	}

	return
}

func aptKeyParseList(list string) (keys []string) {
	keyRe := regexp.MustCompile("^pub.+/(.+) [0-9-]+$")
	for _, line := range strings.Split(list, "\n") {
		v := keyRe.FindStringSubmatch(line)
		if v != nil {
			keys = append(keys, v[1])
		}
	}

	return
}
