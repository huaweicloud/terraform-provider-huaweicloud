---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_backups"
description: |-
  Use this data source to get the list of GaussDB MySQL backups.
---

# huaweicloud_gaussdb_mysql_backups

Use this data source to get the list of GaussDB MySQL backups.

## Example Usage

```hcl
variable "backup_name" {}

data "huaweicloud_gaussdb_mysql_backups" "test" {
  name = var.backup_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the ID of the GaussDB MySQL instance.

* `backup_id` - (Optional, String) Specifies the ID of the backup.

* `name` - (Optional, String) Specifies the backup name.

* `backup_type` - (Optional, String) Specifies the backup type.
  Value options:
  + **auto**: automated full backup.
  + **manual**: manual full backup.

* `instance_name` - (Optional, String) Specifies the instance name.

* `begin_time` - (Optional, String) Specifies the backup start time. The format is **yyyy-mm-ddThh:mm:ssZ**.

* `end_time` - (Optional, String) Specifies the backup end time.The format is **yyyy-mm-ddThh:mm:ssZ**.
  The end time must be later than the start time.

* `status` - (Optional, String) Specifies the backup type.
  Value options:
  + **BUILDING**: The backup is in progress.
  + **COMPLETED**: The backup is complete.
  + **FAILED**: The backup failed.
  + **AVAILABLE**: The backup is available.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `backups` - Indicates the list of backups.

  The [backups](#backups_struct) structure is documented below.

<a name="backups_struct"></a>
The `backups` block supports:

* `id` - Indicates the ID of the backup.

* `name` - Indicates the name of the backup.

* `instance_name` - Indicates the name of the GaussDB MySQL instance.

* `instance_id` - Indicates the ID of the GaussDB MySQL instance.

* `begin_time` - Indicates the backup start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - Indicates the backup end time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `take_up_time` - Indicates the backup duration in minutes.

* `type` - Indicates the backup type.

* `size` - Indicates the backup size in MB.

* `datastore` - Indicates the database information.

  The [datastore](#backups_datastore_struct) structure is documented below.

* `status` - Indicates the backup type.

* `description` - Indicates the description of the backup.

<a name="backups_datastore_struct"></a>
The `datastore` block supports:

* `type` - Indicates the database engine.

* `version` - Indicates the database version.
