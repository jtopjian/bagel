Connections
===========

### Table of Contents

* [Connection Drivers](#connection-drivers)
    * [ssh](#ssh)

Connections are methods of connecting to a target node.

Connections are defined as follows:

```yaml
connections:
  name-of-connection:
    type: connection-driver
    options:
      key: value
      key: value
```

Connection Drivers
------------------

### ssh

The `ssh` driver will connect to a hsot via SSH.

#### example

```yaml
connections:
  name-of-connection:
    type: ssh
    options:
      private_key: /path/to/id_rsa
      port: 22
      shell: /bin/bash
      timeout: 120
      user: root
      agent: true/false
      bastion_host: bastion-host
      bastion_private_key: /path/to/id_rsa
      bastion_user: root
      bastion_port: 22
```

#### options

* `agent` (optional) - Whether or not to use an SSH agent.

* `private_key` (optional) - The SSH private key to connect to the host with.
  If not defined and if `agent` is not `true`, `~/.ssh/id_rsa` will be used.

* `port` (optional) - The port to connect to on the host. Defaults to 22.

* `shell` (optional) - The shell to use on the remote host. Defaults
  to `/bin/bash`.

* `timeout` (optional) - The amount of time (in seconds) to attempt to connect
  to the remote host.

* `user` (optional) - The user to connect to on the remote host. Defaults to
  `root`.

* `bastion_host` (optional) - The bastion host.

* `bastion_user` (optional) - The user to connect to on the bastion host.

* `bastion_private_key` (optional) - The SSH private key to connect to the
  bastion host with. If not defined and if `agent` is not `true`,
  `~/.ssh/id_rsa` will be used.

* `bastion_port` (optional) The port to connect to on the bastion host.
