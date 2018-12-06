cron.Entry
==========

`cron.Entry` will manage a cron entry.

## example

```lua
change, err = cron.Entry({
  name = "foo",
  command = "ls",
  minute = 1,
  hour = 2,
})
```

## options

* `name` (required) - A descriptive name for the cron entry.

* `state` (optional) - The state of the entry. This can either be:
  a version number, `present` or `absent`. Defaults to `present`.

* `command` (required) - The command to perform.

* `minute` (optional) - The minute entry of the cron. Defaults to `*`.

* `hour` (optional) - The hour entry of the cron. Defaults to `*`.

* `day_of_month` (optional) - The day_of_month entry of the cron. Defaults to `*`.

* `month` (optional) - The month entry of the cron. Defaults to `*`.

* `day_of_week` (optional) - The day_of_week entry of the cron. Defaults to `*`.

* `sudo` (optional) - Whether or not sudo is required. Valid values are
  `true` or `false`.

* `timeout` (optional) - How long the command should run before it times out.
