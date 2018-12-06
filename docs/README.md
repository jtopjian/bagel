Bagel Documentation
===================

### Table of Contents

* [Run Mode](#run)
* [Deploy Mode](#deploy)
* [Resources](#resources)
* [`bagel.yaml`](#bagel.yaml)

Bagel can be used in two different styles:

Run
---

Run Mode will apply a configuration to the localhost.

To use Run Mode, create a `.lua` file and then run:

```shell
$ bagel run /path/to/file.lua
```

Deploy
------

Deploy Mode will apply a configuration to a set of remote nodes. To use Deploy
Mode, you need to create a Site File.

### Site File

A site file is where you describe the roles which will be applied to nodes.
A site file is called `site.yaml`. This cannot be changed.

A site file contains 3 parts:

#### Roles

A "role" is a configuration that is applied to a node.

See the [Roles](roles.md) doc for more details.

#### Inventories

An "inventory" is a list of nodes you want to apply a role to. Bagel can have
multiple inventories which can point to different nodes or have overlapping
nodes.

See the [Inventories](inventories.md) doc for more details.

#### Connections

A "connection" is how Bagel will connect to a node.

See the [Connections](connections.md) doc for more details.

Resources
---------

Bagel includes some custom Lua functions which help configure a node in an
idempotent fashion.

See the [resources](resources.md) doc for more details.

bagel.yaml
----------

Bagel can be configured with a `bagel.yaml` file. This file can exist in either:

1. `/opt/bagel`
2. Your current directory.

You can set the following in `bagel.yaml`:

* `site_dir`: Where Bagel can find the `site.yaml`. By default, this is `/opt/bagel`.
* `debug`: Whether to enable debugging. By default, this is false.
