log.Info("Starting apt.Key tests")

-- Run apt-get update
_, err = exec.Run({
  cmd = "apt-get update -qq",
})

util.StopIfError("Command failed", err)

-- Install required packages for test
packages = {
  "ca-certificates",
  "gnupg",
}

for i, pkg in ipairs(packages) do
  change, err = apt.Package({
    name = pkg,
  })

  util.StopIfError("Error installing" .. pkg, err)
  assert(change)
end

-- Install rabbit key
change, err = apt.Key({
  name = "6026DFCA",
  remote_key_file = "https://github.com/rabbitmq/signing-keys/releases/download/2.0/rabbitmq-release-signing-key.asc",
})

util.StopIfError("Error installing key", err)
assert(change)

-- Confirm it was installed
x = os.execute("apt-key adv --list-public-keys | grep -q 6026DFCA")
assert(x == 0)

-- Install it again
-- make sure there was no change
change, err = apt.Key({
  name = "6026DFCA",
  remote_key_file = "https://github.com/rabbitmq/signing-keys/releases/download/2.0/rabbitmq-release-signing-key.asc",
})

util.StopIfError("Error installing key", err)
assert(not change)

-- Remove key
change, err = apt.Key({
  name = "6026DFCA",
  remote_key_file = "https://github.com/rabbitmq/signing-keys/releases/download/2.0/rabbitmq-release-signing-key.asc",
  state = "absent",
})

util.StopIfError("Error removing key", err)
assert(change)

-- Confirm it was removed
x = os.execute("apt-key adv --list-public-keys | grep -q 6026DFCA")
assert(x == 1)

-- Remove it again
-- make sure there was no change
change, err = apt.Key({
  name = "6026DFCA",
  remote_key_file = "https://github.com/rabbitmq/signing-keys/releases/download/2.0/rabbitmq-release-signing-key.asc",
  state = "absent",
})

util.StopIfError("Error removing key", err)
assert(not change)
