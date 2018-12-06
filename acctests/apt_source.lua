log.Info("Starting apt.Source tests")

-- Install source file
change, err = apt.Source({
  name         = "rabbitmq",
  uri          = "https://dl.bintray.com/rabbitmq/debian",
  distribution = "bionic",
  component    = "main",
})

util.StopIfError(err, "Error installing apt source file")
assert(change)

-- Confirm it was installed
x = os.execute("grep -q 'https://dl.bintray.com/rabbitmq/debian bionic main' /etc/apt/sources.list.d/rabbitmq.list")
assert(x == 0)

-- Install it again
-- make sure there was no change
change, err = apt.Source({
  name         = "rabbitmq",
  uri          = "https://dl.bintray.com/rabbitmq/debian",
  distribution = "bionic",
  component    = "main",
})

util.StopIfError(err, "Error installing apt source file")
assert(not change)

-- Remove source file
change, err = apt.Source({
  name         = "rabbitmq",
  uri          = "https://dl.bintray.com/rabbitmq/debian",
  distribution = "bionic",
  component    = "main",
  state        = "absent",
})

-- Confirm it was deleted
x = os.execute("ls /etc/apt/sources.list.d/rabbitmq.list")
assert(x == 1)
