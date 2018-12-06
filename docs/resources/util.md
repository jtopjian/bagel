Utility Functions
=================

util.LogIfError
---------------

If an error occurred, log it, but continue applying the role.

### Example

```lua
change, err = apt.Package({
  name = "sl",
})

util.LogIfError("Error installing sl", err)
```


util.StopIfError
----------------

If an error occurred, log it, and halt applying the role.

### Example

```lua
change, err = apt.Package({
  name = "sl",
})

util.StopIfError("Error installing sl", err)
```
