log.Info("Starting cron.Entry tests")

-- Run apt-get update
_, err = exec.Run({
  cmd = "apt-get update -qq",
})

util.StopIfError("Command failed", err)

-- Install required packages for test
packages = {
  "cron",
}

for i, pkg in ipairs(packages) do
  change, err = apt.Package({
    name = pkg,
  })

  util.StopIfError("Error installing" .. pkg, err)
  assert(change)
end

-- Add cron entry
change, err = cron.Entry({
  name = "foo",
  command = "ls",
  minute = 1,
  hour = 2,
})

util.StopIfError("Error creating cron entry", err)
assert(change)

-- Confirm it was installed
x = os.execute("grep '1 2 \\* \\* \\* ls # foo' /var/spool/cron/crontabs/root")
assert(x == 0)

-- Add it again
-- verify no changes happened
change, err = cron.Entry({
  name = "foo",
  command = "ls",
  minute = 1,
  hour = 2,
})

util.StopIfError("Error creating cron entry", err)
assert(not change)

-- Remove the cron entry
change, err = cron.Entry({
  name = "foo",
  command = "ls",
  minute = 1,
  hour = 2,
  state = "absent",
})

util.StopIfError("Error creating cron entry", err)
assert(change)

-- Confirm it was removed
x = os.execute("grep '1 2 \\* \\* \\* ls # foo' /var/spool/cron/crontabs/root")
assert(x == 1)
