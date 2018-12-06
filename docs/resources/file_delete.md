file.Delete
===========

`file.Delete` will delete a file if it exists.

## example

```lua
change, err = file.Delete({
  path = "/path/to/file",
})
```

## options

* `path` (required) - The path to the file to delete.

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
