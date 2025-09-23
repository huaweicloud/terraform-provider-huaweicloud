---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_hba"
description: ""
---

# huaweicloud_rds_pg_hba

Manages an RDS PostgreSQL hba resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_pg_hba" "test" {
  instance_id = var.instance_id

  host_based_authentications {
    type     = "host"
    database = "all"
    user     = "all"
    address  = "0.0.0.0/0"
    method   = "md5"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS PostgreSQL instance.

* `host_based_authentications` - (Required, List) Specifies the list of host based authentications.
The [host_based_authentications](#PgHba_HostBasedAuthentication) structure is documented below.

<a name="PgHba_HostBasedAuthentication"></a>
The `host_based_authentications` block supports:

* `type` - (Required, String) Specifies the connection type. Value options: **host**, **hostssl** and **hostnossl**.

* `database` - (Required, String) Specifies the database name other than **template0** and **template1**.
  + **all** indicates all databases of the DB instance.
  + Use commas (,) to separate multiple databases.

* `user` - (Required, String) Specifies the name of a user other than **rdsAdmin**, **rdsMetric**, **rdsBackup**,
  **rdsRepl** and **rdsProxy**.
  + **all** indicates all database users of the DB instance.
  + Use commas (,) to separate multiple user names.

* `address` - (Required, String) Specifies the client IP address.
  + **0.0.0.0/0** indicates that the user can access the database from any IP address.

* `method` - (Required, String) Specifies the authentication mode. Value options: **reject**, **md5** and
  **scram-sha-256**.

* `mask` - (Optional, String) Specifies the subnet mask. It is mandatory when `address` does not contain mask.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the instance ID.

## Import

The rds PostgreSQL hba can be imported using the `instance_id`, e.g.

```bash
$ terraform import huaweicloud_rds_pg_hba.test <instance_id>
```
