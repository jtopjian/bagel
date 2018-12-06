file.Push
=========

`file.Push` will push a file to a remote node.

## example

```lua
info, err = file.Push({
  source = "/path/to/local/file"
  destination = "/path/to/remote/file"
})

```

## options

* `source` (required) - The path to the source file on the local node.

* `destination` (required) - The path to the destination file on the remote node.

* `timeout` (optional) - How long the command should run before it times out.

## returns

* `applied` - Whether a change was happened.

* `exists` - Whether or not the file exists.

* `info` - Information about the file.

* `success` - If the action was successful.

* `timeout` - If a timeout happened.

### The `info` table contains

* `name` - The name of the file.

* `uid` - The UID of the file.

* `gid` - The UID of the file.

* `type` - The type of the file.

* `size` - The size of the file.

* `mode` - The mode of the file.
