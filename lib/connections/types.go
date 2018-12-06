package connections

import (
	"io"
	"os"

	"github.com/yuin/gopher-lua"
)

// Connection is an interface which specifies what drivers
// must implement.
type Connection interface {
	Connect() error
	Close()

	RunCommand(RunOpts) (*RunResult, error)

	FileInfo(FileOpts) (*FileResult, error)
	FileDelete(FileOpts) (*FileResult, error)
	FileUpload(CopyFileOpts) (*FileResult, error)
	FileDownload(CopyFileOpts) (*FileResult, error)
}

// RunOpts represents options for running commands.
type RunOpts struct {
	Command string
	Timeout int
	Log     *io.Writer
}

// RunResult respresents the result of an command execution.
type RunResult struct {
	ExitCode int
	Stderr   string
	Stdout   string
	Timeout  bool
	Applied  bool
}

// ToLTable converts a RunResult to a GopherLua table.
func (r RunResult) ToLTable(L *lua.LState) *lua.LTable {
	ret := L.NewTable()
	ret.RawSetString("exit_code", lua.LNumber(r.ExitCode))
	if r.Stderr == "" {
		ret.RawSetString("stderr", lua.LNil)
	} else {
		ret.RawSetString("stderr", lua.LString(r.Stderr))
	}
	ret.RawSetString("stdout", lua.LString(r.Stdout))
	ret.RawSetString("timeout", lua.LBool(r.Timeout))
	ret.RawSetString("applied", lua.LBool(r.Applied))

	return ret
}

// CopyFileOpts represents options for copying files.
type CopyFileOpts struct {
	Source      string
	Destination string
	UID         int
	GID         int
	Mode        os.FileMode
	Timeout     int
}

// FileOpts represents options for managing a generic file.
type FileOpts struct {
	Path    string
	UID     int
	GID     int
	Mode    os.FileMode
	Timeout int
}

// FileResult represents the result of an file action.
type FileResult struct {
	Exists   bool
	Success  bool
	Timeout  bool
	Applied  bool
	FileInfo FileInfo
}

// ToLTable converts a FileResult to a GopherLua table.
func (r FileResult) ToLTable(L *lua.LState) *lua.LTable {
	ret := L.NewTable()
	ret.RawSetString("exists", lua.LBool(r.Exists))
	ret.RawSetString("success", lua.LBool(r.Success))
	ret.RawSetString("timeout", lua.LBool(r.Timeout))
	ret.RawSetString("applied", lua.LBool(r.Applied))
	ret.RawSetString("info", r.FileInfo.ToLTable(L))

	return ret
}

// FileInfo represents information about a file.
type FileInfo struct {
	Name string
	UID  int
	GID  int
	Type string
	Size int64
	Mode int
}

func (r FileInfo) ToLTable(L *lua.LState) *lua.LTable {
	ret := L.NewTable()

	ret.RawSetString("name", lua.LString(r.Name))
	ret.RawSetString("uid", lua.LNumber(r.UID))
	ret.RawSetString("gid", lua.LNumber(r.GID))
	ret.RawSetString("type", lua.LString(r.Type))
	ret.RawSetString("size", lua.LNumber(r.Size))
	ret.RawSetString("mode", lua.LNumber(r.Mode))

	return ret
}
