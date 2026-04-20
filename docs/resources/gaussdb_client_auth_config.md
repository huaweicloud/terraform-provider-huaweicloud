---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_client_auth_config"
description: |-
  Manages a GaussDB client auth config resource within HuaweiCloud.
---

# huaweicloud_gaussdb_client_auth_config

Manages a GaussDB client auth config resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_client_auth_config" "test" {
  instance_id = var.instance_id
  type        = "host"
  database    = "all"
  user        = "root"
  address     = "10.10.0.0/16"
  method      = "reject"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GaussDB instance.

* `type` - (Required, String, NonUpdatable) Specifies the client connection type. Valid values include **host**, **hostssl**, **hostnossl**.

* `database` - (Required, String, NonUpdatable) Specifies the name of the database that the record matches.
  The value can be **all** or an existing database name.

* `user` - (Required, String, NonUpdatable) Specifies the name of the database user that the record matches.
  The value can be **all** or an existing username.

* `address` - (Required, String, NonUpdatable) Specifies the IP address range that the record matches.
  The value must be in CIDR format (e.g., `10.10.0.0/16`).

* `method` - (Required, String) Specifies the authentication method used for the connection.
  Valid values include **md5**, **sha256**, **sm3**, **reject**, **cert**, etc.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is formatted as `<instance_id>:<type>:<database>:<user>:<address>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 90 minutes.
* `update` - Default is 90 minutes.
* `delete` - Default is 90 minutes.

## Import

The GaussDB client auth config can be imported using the `instance_id`, `type`, `database`, `user` and `address`
separated by colons, e.g.

```bash
$ terraform import huaweicloud_gaussdb_client_auth_config.test <instance_id>:<type>:<database>:<user>:<address>
```
