Inventories
===========

### Table of Contents

* [Inventory Drivers](#inventory-drivers)
    * [textfile](#textfile)

Inventories define the ndoes which a configuration will be applied to.
Inventories can be dynamically discovered or statically defined.

Inventories are defined as follows:

```yaml
inventories:
  name-of-inventory:
    type: inventory-driver
    options:
      key: value
      key: value
    connection: connection-driver
```

Inventory Drivers
-----------------

Bagel currenty supports the following Inventory Drivers:

### textfile

The `textfile` driver will read nodes defined in a plain text file.

#### example

```yaml
inventories:
  my_nodes:
    type: textfile
    options:
      file: /path/to/nodes.txt
```

#### options

* `file` (required) - The text file which defines the hosts. Each
line of the text file must contain only the resolvable name or IP
address of the host. An example file is:

```
# comment
host1.example.com
host2.example.com
// host3.example.com
192.168.100.1
fe80::f816:3eff:fe8c:c73a
```
