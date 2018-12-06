**NOTE:** Until there's an actual binary release made in [releases](/releases), consider this
a proof of concept and not guaranteed to work at all.

Bagel
=====

Yet another configuration management tool I dreamed up.

Bagel builds off of previous projects such as [Waffles](https://github.com/wffls/waffles)
and [Yak](https://github.com/jtopjian/yak). It supports:

* The ability to write configuration management manifests in Lua. This means
  you can use `if` conditionals and `for` loops natively.

* Being able to work with remote communication protocols other than SSH. This
  means Bagel can (and probably will) support protocols such as LXD, Docker, etc.
  This can extend beyond simply remotely executing a command: when available,
  Bagel can take advantage of remote file API rather than just echo'ing or
  cat'ing content to a file.

* Is written in Go and distributed as a single Go binary for ease of use.
  A plugin system might be available in the future since it might be quite
  useful.

Quickstart
----------

1. Download Bagel.
2. Create a directory:

  ```bash
  $ mkdir /opt/bagel
  $ cd /opt/bagel
  ```

3. Add some hsots to a file:

  ```bash
  $ echo example1.com >> hosts.txt
  $ echo example2.com >> hosts.txt
  ```

4. Create a `site.yaml` file::

  ```yaml
  roles:
    hello:
      inventories:
        - my_hosts

  inventories:
    my_hosts:
      type: textfile
      options:
        file: /opt/bagel/hosts.txt
      connection: ssh

  connections:
    ssh:
      type: ssh
  ```

5. Create a "hello" role in `/opt/bagel/roles/hello.lua`:

  ```lua
  log.Info("Hello, World!")

  change, err = apt.Package({
      name = "sl",
  })

  util.StopIfErr("Unable to install sl", err)
  ```

6. Run Bagel:

  ```shell
  $ bagel deploy
  ```
Documentation
-------------

See the [docs](/docs) directory.

You can also check out the [acceptance tests](/acctests) for examples.

Why??
-----

Same as always: to scratch an itch and create something that didn't exist before.

Building from Source
--------------------

```bash
$ go get -u github.com/jtopjian/bagel/...
$ cd $GOPATH/src/github.com/jtopjian/bagel
$ make build
# or
$ make install
```
