---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_backups"
description: |-
  Use this data source to get the list of GaussDB OpenGauss backups.
---

# huaweicloud_gaussdb_opengauss_backups

Use this data source to get the list of GaussDB OpenGauss backups.

## Example Usage

```hcl
variable "instance_id" {}
variable "backup_id" {}

data "huaweicloud_gaussdb_opengauss_backups" "test" {
  instance_id = var.instance_id
  backup_id   = var.backup_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the ID of the GaussDB OpenGauss instance

* `backup_id` - (Optional, String) Specifies the ID of the backup.

* `backup_type` - (Optional, String) Specifies the backup type.
  Value options:
  + **auto**: instance-level automated full backup.
  + **manual**: instance-level manual full backup.

* `begin_time` - (Optional, String) Specifies the query start time in the **yyyy-mm-ddThh:mm:ssZ** format.
  It can be used together with `end_time`. If `end_time` is not used, the backups created after begin_time are
  returned. If `end_time` is used, the backups created between `begin_time` and `end_time` are returned.

* `end_time` - (Optional, String) Specifies the query end time in the **yyyy-mm-ddThh:mm:ssZ** format.
  It must be later than the start time. It can be used together with `begin_time`. If `begin_time` is not used, the
  backups created before `end_time` are returned. If `begin_time` is used, the backups created between
  `begin_time` and `end_time` are returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `backups` - Indicates the list of backups.

  The [backups](#backups_struct) structure is documented below.

<a name="backups_struct"></a>
The `backups` block supports:

* `id` - Indicates the ID of the backup.

* `name` - Indicates the name of the backup.

* `instance_id` - Indicates the ID of the GaussDB OpenGauss instance.

* `description` - Indicates the description of the backup.

* `begin_time` - Indicates the backup start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - Indicates the backup end time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `type` - Indicates the backup type.

* `size` - Indicates the backup size in MB.

* `status` - Indicates the backup status.
  The value can be:
  + **BUILDING**: Backup in progress
  + **COMPLETED**: Backup completed
  + **FAILED**: Backup failed

* `datastore` - Indicates the database information.

  The [datastore](#backups_datastore_struct) structure is documented below.

<a name="backups_datastore_struct"></a>
The `datastore` block supports:

* `type` - Indicates the DB engine.
  The value is case-insensitive and can be: **GaussDB**.

* `version` - Indicates the DB engine version
