apt.PPA
=======

`apt.PPA` will manage an apt PPA repository.

## example

```lua
change, err = apt.PPA({
  name = "chris-lea/redis-server",
})
```

## options

* `name` (required) - The short ID of the key.

* `state` (optional) - The state of the key. This can either be:
  `present` or `absent`. Defaults to `present`.

* `refresh` (optional) - Whether to perform an `apt-get update`
  when the state changes. Defaults to `false`.

* `sudo` (optional) - Whether or not sudo is required. Valid values are
  `true` or `false`.

* `timeout` (optional) - How long the command should run before it times out.
