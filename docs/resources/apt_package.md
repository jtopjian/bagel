apt.Package
===========

`apt.Package` will manage an apt package.

## example

```lua
change, err = apt.Package({
  name = "sl",
})
```

## options

* `name` (required) - The short ID of the key.

* `state` (optional) - The state of the key. This can either be:
  `present` or `absent`. Defaults to `present`.

* `sudo` (optional) - Whether or not sudo is required. Valid values are
  `true` or `false`.

* `timeout` (optional) - How long the command should run before it times out.
