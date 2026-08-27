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
data "huaweicloud_dds_backups" "test" {}
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

* `instance_id` - (Optional, String) Specifies the ID of the DDS instance from which the backup was created.

* `begin_time` - (Optional, String) Specifies the start time of the query. The format is **yyyy-mm-dd hh:mm:ss**.
  The value is in UTC format. It's required with `end_time`.

* `end_time` - (Optional, String) Specifies the end time of the query. The format is **yyyy-mm-dd hh:mm:ss**.
  The value is in UTC format. It's required with `begin_time`.

* `mode` - (Optional, String) Specifies the DDS instance mode.
  The valid values are **Sharding**, **ReplicaSet** and **Single**.

* `instance_name` - (Optional, String) Specifies the name of the DDS instance for which the backup is created.
  Fuzzy match is supported.

* `backup_name` - (Optional, String) Specifies the backup name. Fuzzy match is supported.

* `status` - (Optional, String) Specifies the backup status. Valid values are:
  + **BUILDING**: Backup in progress.
  + **COMPLETED**: Backup completed.
  + **FAILED**: Backup failed.

* `description` - (Optional, String) Specifies the backup description. Fuzzy match is supported.

* `order_field` - (Optional, String) Specifies the sort field. It's required with `order_rule`.
  The valid values are as follows:
  + **name**: Indicates backup name.
  + **instanceName**: Indicates instance name.
  + **type**: Indicates backup type.
  + **datastoreType**: Indicates engine type.
  + **beginTime**: Indicates start time.
  + **status**: Indicates backup status.

  If this parameter is not specified, backups are sorted in descending order based on the backup start time.

* `order_rule` - (Optional, String) Specifies the sort rule. It's required with `order_field`.
  The valid values are as follows:
  + **asc**: Indicates ascending order.
  + **desc**: Indicates descending order.

  If this parameter is not specified, backups are sorted in descending order based on the backup start time.

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

* `instance_status` - Indicates the instance status.
  + **normal**: An instance is running normally.
  + **abnormal**: An instance is abnormal.
  + **creating**: An instance is being created.
  + **frozen**: An instance is frozen.
  + **data_disk_full**: The storage space is full.
  + **enlargefail**: Nodes failed to be added to the instance.

* `instance_mode` - Indicates the instance mode..

* `is_instance_restoring` - Whether the current instance is being restored or checked.

* `backup_method` - Indicates the backup method.
  + **Snapshot**: Indicates snapshot backup.
  + **Physical**: Indicates physical backup.
  + **Logical**: Indicates logical backup.
  + **Incremental**: Indicates incremental backup.

* `kms_enable` - Whether KMS encryption is enabled.

* `deletable` - Whether the backup can be deleted.

<a name="backups_datastore_struct"></a>
The `datastore` block supports:

* `type` - Indicates the DB engine.

* `version` - Indicates the database version.
