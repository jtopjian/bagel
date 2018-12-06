package testing

import (
	"testing"

	"github.com/jtopjian/bagel/lib/connections"
	"github.com/stretchr/testify/assert"
)

func TestSSH_Basic(t *testing.T) {
	options := map[string]interface{}{
		"host":        "localhost",
		"user":        "ubuntu",
		"private_key": "/root/.ssh/id_rsa",
		"shell":       "/bin/bash",
		"timeout":     5,
	}

	ssh, err := connections.New("ssh", options)
	if err != nil {
		t.Fatal(err)
	}

	if err := ssh.Connect(); err != nil {
		t.Fatal(err)
	}

	ro := connections.RunOpts{
		Command: "echo hi",
	}

	rr, err := ssh.RunCommand(ro)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "hi", rr.Stdout)

	ro.Command = "asdf"
	rr, err = ssh.RunCommand(ro)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "/bin/bash: asdf: command not found", rr.Stderr)

	ro.Command = `foo=bar; sleep 1; echo foobar >&2; echo \$foo ; echo 123 >&2`
	rr, err = ssh.RunCommand(ro)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "bar", rr.Stdout)
	assert.Equal(t, "foobar\n123", rr.Stderr)
}

func TestSSH_CommandTimeout(t *testing.T) {
	options := map[string]interface{}{
		"host":        "localhost",
		"user":        "ubuntu",
		"private_key": "/root/.ssh/id_rsa",
		"shell":       "/bin/bash",
	}

	ssh, err := connections.New("ssh", options)
	if err != nil {
		t.Fatal(err)
	}

	if err := ssh.Connect(); err != nil {
		t.Fatal(err)
	}

	ro := connections.RunOpts{
		Command: "sleep 6; echo timeout",
		Timeout: 5,
	}

	rr, err := ssh.RunCommand(ro)
	assert.Equal(t, true, rr.Timeout)
}

func TestSSH_ConnectTimeout(t *testing.T) {
	options := map[string]interface{}{
		"host":        "localhost2",
		"user":        "ubuntu",
		"private_key": "/root/.ssh/id_rsa",
		"shell":       "/bin/bash",
		"timeout":     5,
	}

	ssh, err := connections.New("ssh", options)
	if err != nil {
		t.Fatal(err)
	}

	err = ssh.Connect()
	assert.Equal(t, "timed out connecting to localhost2:22", err.Error())
}

func TestSSH_Bastion(t *testing.T) {

	options := map[string]interface{}{
		"host":                "localhost",
		"user":                "ubuntu",
		"private_key":         "/root/.ssh/id_rsa",
		"shell":               "/bin/bash",
		"timeout":             5,
		"bastion_host":        "localhost",
		"bastion_user":        "ubuntu",
		"bastion_private_key": "/root/.ssh/id_rsa",
	}

	ssh, err := connections.New("ssh", options)
	if err != nil {
		t.Fatal(err)
	}

	if err := ssh.Connect(); err != nil {
		t.Fatal(err)
	}

	ro := connections.RunOpts{
		Command: "echo hi",
	}

	rr, err := ssh.RunCommand(ro)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "hi", rr.Stdout)
}

func TestSSH_CopyFileDelete(t *testing.T) {
	options := map[string]interface{}{
		"host":        "localhost",
		"user":        "ubuntu",
		"private_key": "/root/.ssh/id_rsa",
		"shell":       "/bin/bash",
	}

	ssh, err := connections.New("ssh", options)
	if err != nil {
		t.Fatal(err)
	}

	if err := ssh.Connect(); err != nil {
		t.Fatal(err)
	}

	cfo := connections.CopyFileOpts{
		Source:      "fixtures/hello.txt",
		Destination: "/tmp/bagelfoo.txt",
	}

	fr, err := ssh.FileUpload(cfo)
	if err != nil {
		t.Fatal(err)
	}

	ro := connections.RunOpts{
		Command: "cat /tmp/bagelfoo.txt",
	}

	rr, err := ssh.RunCommand(ro)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, true, fr.Success)
	assert.Equal(t, "Hello, World!", rr.Stdout)

	fo := connections.FileOpts{
		Path: "/tmp/bagelfoo.txt",
	}

	fr, err = ssh.FileDelete(fo)
	if err != nil {
		t.Fatal(err)
	}

	ro.Command = "stat /tmp/bagelfoo.txt"
	rr, err = ssh.RunCommand(ro)
	if err != nil {
		t.Fatal(err)
	}

	if rr.ExitCode != 1 {
		t.Fatalf("file still exists")
	}

}
