Roles
=====

Roles are a set of configurations that are applied to nodes. A role is just a
Lua script with some Bagel-specific functions built-in.

Role are stored in `/opt/bagel/roles`, for example, `/opt/bagel/roles/memcached.lua`.

Roles are defined in the `/opt/bagel/site.yaml` file like so:

```yaml
roles:
  name-of-role:
    inventories
      - inventory_1
      - inventory_2
```

## Options

* `inventories` (Required) - The inventories to apply the role to.
