---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_backups"
description: |-
  Use this data source to get the list of DDS instance backups.
---

# huaweicloud_dds_backups

Use this data source to get the list of DDS instance backups.

## Example Usage

```hcl
data "huaweicloud_dds_backups" "all" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `backup_id` - (Optional, String) Specifies the backup ID.
  If the backup ID belongs to an automated incremental backup, the `instance_id` is required.

* `backup_type` - (Optional, String) Specifies the backup type. Valid values are:
  + **Auto**: Indicates automated full backup.
  + **Manual**: Indicates manual full backup.
  + **Incremental**: Indicates automated incremental backup.

* `instance_id` - (Optional, String) Specifies the ID of the DB instance from which the backup was created.

* `begin_time` - (Optional, String) Specifies the start time of the query. The format is **yyyy-mm-dd hh:mm:ss**.
  The value is in UTC format. It's required with `end_time`.

* `end_time` - (Optional, String) Specifies the end time of the query. The format is **yyyy-mm-dd hh:mm:ss**.
  The value is in UTC format. It's required with `begin_time`.

* `mode` - (Optional, String) Specifies the DB instance mode. Valid values are **Sharding** and **ReplicaSet**.

* `instance_name` - (Optional, String) Specifies the name of the DB instance for which the backup is created.

* `backup_name` - (Optional, String) Specifies the backup name.

* `status` - (Optional, String) Specifies the backup status. Valid values are:
  + **BUILDING**: Backup in progress.
  + **COMPLETED**: Backup completed.
  + **FAILED**: Backup failed.
  + **DISABLED**: Backup being deleted.

* `description` - (Optional, String) Specifies the backup description.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `backups` - Indicates the backup list.
  The [backups](#backups_struct) structure is documented below.

<a name="backups_struct"></a>
The `backups` block supports:

* `id` - Indicates the backup ID.

* `name` - Indicates the backup name.

* `instance_id` - Indicates the ID of the DB instance from which the backup was created.

* `instance_name` - Indicates the name of the DB instance for which the backup is created.

* `type` - Indicates the backup type.

* `size` - Indicates the backup size in KB.

* `datastore` - Indicates the database version.
  The [datastore](#backups_datastore_struct) structure is documented below.

* `begin_time` - Indicates the backup start time.

* `end_time` - Indicates the backup end time.

* `status` - Indicates the backup status.

* `description` - Indicates the backup description.

<a name="backups_datastore_struct"></a>
The `datastore` block supports:

* `type` - Indicates the DB engine.

* `version` - Indicates the database version.
