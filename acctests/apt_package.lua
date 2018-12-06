log.Info("Starting apt.Package tests")

-- Run apt-get update
_, err = exec.Run({
  cmd = "apt-get update -qq",
})

util.StopIfError("Command failed", err)

-- Install sl
change, err = apt.Package({
  name = "sl",
})

util.StopIfError("Error installing sl", err)
assert(change)

-- Confirm it was installed
x = os.execute("dpkg -l | grep -q ' sl '")
assert(x == 0)

-- Install it again
-- make sure there was no change
change, err = apt.Package({
  name = "sl",
})

util.StopIfError("Error installing sl", err)
assert(not change)

-- Remove sl
change, err = apt.Package({
  name = "sl",
  state = "absent",
})

util.StopIfError("Error removing sl", err)
assert(change)

-- Confirm it was removed
x = os.execute("dpkg -l | grep -q ' sl '")
assert(x == 1)

-- Remove it again
-- make sure there was no change
change, err = apt.Package({
  name = "sl",
  state = "absent",
})

util.StopIfError("Error removing sl", err)
assert(not change)

-- Install a specific version of sl
change, err = apt.Package({
  name = "sl",
  state = "3.03-17build2",
})

util.StopIfError("Error installing sl", err)
assert(change)

-- Confirm it was installed
x = os.execute("dpkg -l | grep ' sl '")
assert(x == 0)

-- Install a specific version again
-- make sure there was no change
change, err = apt.Package({
  name = "sl",
  state = "3.03-17build2",
})

util.StopIfError("Error installing sl", err)
assert(not change)

-- Install a non-specific version
-- make sure there was no change
change, err = apt.Package({
  name = "sl",
})

util.StopIfError("Error installing sl", err)
assert(not change)

-- Remove sl
change, err = apt.Package({
  name = "sl",
  state = "absent",
})

util.StopIfError("Error removing sl", err)
assert(change)
