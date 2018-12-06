log.Info("Starting apt.PPA tests")

-- Run apt-get update
_, err = exec.Run({
  cmd = "apt-get update -qq",
})

util.StopIfError("Command failed", err)

-- Install required packages for test
packages = {
  "software-properties-common",
}

for i, pkg in ipairs(packages) do
  change, err = apt.Package({
    name = pkg,
  })

  util.StopIfError("Error installing" .. pkg, err)
  assert(change)
end

-- Install redis PPA
change, err = apt.PPA({
  name = "chris-lea/redis-server",
})

util.StopIfError("Error installing PPA", err)
assert(change)

-- Confirm it was installed
x = os.execute("ls /etc/apt/sources.list.d/")
x = os.execute("ls /etc/apt/sources.list.d/chris-lea-ubuntu-redis-server-bionic.list")
assert(x == 0)

-- Install it again
-- make sure there was no change
change, err = apt.PPA({
  name = "chris-lea/redis-server",
})

util.StopIfError("Error installing PPA", err)
assert(not change)
