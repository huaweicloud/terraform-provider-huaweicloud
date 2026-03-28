---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_remote_databases"
description: |-
  Use this data source to get the remote databases when subscribing to remote SQL Server instances.
---

# huaweicloud_rds_remote_databases

Use this data source to get the remote databases when subscribing to remote SQL Server instances.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_remote_databases" "test" {
  instance_id         = var.instance_id
  server_ip           = "192.168.0.230"
  server_port         = "1433"
  login_user_name     = "rdsuser"
  login_user_password = "test_1234"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `server_ip` - (Required, String) Specifies the IP address of the remote database.

* `server_port` - (Required, String) Specifies the port of the remote database.

* `login_user_name` - (Required, String) Specifies the user for logging in to the remote database.

* `login_user_password` - (Required, String) Specifies the password for the database user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `databases` - Indicates the list of remote databases.
  The [databases](#databases_struct) structure is documented below.

<a name="databases_struct"></a>
The `databases` block supports:

* `name` - Indicates the database name.

* `character_set` - Indicates the character set used by the database.

* `state` - Indicates the database status. The value can be:
  + **Creating**: The database is being created.
  + **Running**: The database is running.
  + **Deleting**: The database is being deleted.
  + **NotExists**: The database does not exist.
