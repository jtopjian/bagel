file.Exists
===========

`file.Exists` will verify if a file exists and return information.

## example

```lua
info, err = file.Exists({
  path = "/path/to/file",
})
```

## options

* `path` (required) - The path to the file to delete.

* `timeout` (optional) - How long the command should run before it times out.

## returns

* `exists` - Whether or not the file exists.

* `success` - If the check was successful.

* `timeout` - If a timeout happened.

* `info` - Information about the file.

### The `info` table contains

* `name` - The name of the file.

* `uid` - The UID of the file.

* `gid` - The UID of the file.

* `type` - The type of the file.

* `size` - The size of the file.

* `mode` - The mode of the file.
