apt.Key
=======

`apt.Key` will manage an apt key.

## example

```lua
change, err = apt.Key({
  name = "6026DFCA",
  remote_key_file = "https://github.com/rabbitmq/signing-keys/releases/download/2.0/rabbitmq-release-signing-key.asc",
})
```

## options

* `name` (required) - The short ID of the key.

* `state` (optional) - The state of the key. This can either be:
  `present` or `absent`. Defaults to `present`.

* `remote_key_file` (optional) - The URL to a public key. Cannot be
  used with `key_server`.

* `key_server` (optional) - The remote server to obtain the key from.
  Cannot be used with `remote_key_file`.

* `sudo` (optional) - Whether or not sudo is required. Valid values are
  `true` or `false`.

* `timeout` (optional) - How long the command should run before it times out.
