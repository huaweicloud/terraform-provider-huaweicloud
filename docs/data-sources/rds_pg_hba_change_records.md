---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_hba_change_records"
description: |-
  Use this data source to get the pg_hba.conf change history of a DB instance.
---

# huaweicloud_rds_pg_hba_change_records

Use this data source to get the pg_hba.conf change history of a DB instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_pg_hba_change_records" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `start_time` - (Optional, String) Specifies the start time. If it is not specified, **00:00 (UTC time zone)** on the
  current day is used by default.

* `end_time` - (Optional, String) Specifies the end time. If it is not specified, the current time (UTC time zone) is
  used by default.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `pg_hba_change_records` - Indicates the list of pg_hba.conf change history.
  The [pg_hba_change_records](#pg_hba_change_records_struct) structure is documented below.

<a name="pg_hba_change_records_struct"></a>
The `pg_hba_change_records` block supports:

* `status` - Indicates the change result. The value can be:
  + **success**: The change has taken effect.
  + **failed**: The change did not take effect.
  + **setting**: The change is in progress.

* `time` - Indicates the time when the change was made.

* `fail_reason` - Indicates the reason for a change failure.

* `before_confs` - Indicates the original values.
  The [before_confs](#before_confs_struct) structure is documented below.

* `after_confs` - Indicates the new values.
  The [after_confs](#after_confs_struct) structure is documented below.

<a name="before_confs_struct"></a>
The `before_confs` block supports:

* `type` - Indicates the connection type.

* `database` - Indicates the database name.

* `user` - Indicates the name of a user.

* `address` - Indicates the client IP address.

* `mask` - Indicates the subnet mask.

* `method` - Indicates the authentication mode.

* `priority` - Indicates the configuration priority.

<a name="after_confs_struct"></a>
The `after_confs` block supports:

* `type` - Indicates the connection type.

* `database` - Indicates the database name.

* `user` - Indicates the name of a user.

* `address` - Indicates the client IP address.

* `mask` - Indicates the subnet mask.

* `method` - Indicates the authentication mode.

* `priority` - Indicates the configuration priority.
