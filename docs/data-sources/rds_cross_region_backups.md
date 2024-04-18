---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_cross_region_backups"
description: |-
  Use this data source to get the list of RDS cross-region backups.
---

# huaweicloud_rds_cross_region_backups

Use this data source to get the list of RDS cross-region backups.

## Example Usage

```hcl
variable "instance_id" {}
variable "backup_type" {}

data "huaweicloud_rds_cross_region_backups" "test" {
  instance_id = var.instance_id
  backup_type = var.backup_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the cross-region backup.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `backup_type` - (Required, String) Specifies the type of the cross-region backup.
  Value options:
  + **auto**: automated full backup. Microsoft SQL Server only supports the query of this backup type.
  + **incremental**: automated incremental backup.

* `status` - (Optional, String) Specifies the status of the cross-region backup.
  Value options:
  + **BUILDING**: Backup in progress
  + **COMPLETED**: Backup completed
  + **FAILED**: Backup failed
  + **DELETING**: Backup being deleted

* `backup_id` - (Optional, String) Specifies the ID of the cross-region backup.

* `begin_time` - (Optional, String) Specifies the start time for obtaining the cross-region backup list.
  The format is **yyyy-mm-ddThh:mm:ssZ**. This parameter must be used together with `end_time`.

* `end_time` - (Optional, String) Specifies the end time for obtaining the cross-region backup list.
  The format is **yyyy-mm-ddThh:mm:ssZ**. The end time must be later than the start time.
  This parameter must be used together with `begin_time`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `backups` - Indicates the list of the cross-region backups.

  The [backups](#backups_struct) structure is documented below.

<a name="backups_struct"></a>
The `backups` block supports:

* `id` - Indicates the ID of the cross-region backup.

* `name` - Indicates the name of the cross-region backup.

* `type` - Indicates the type of the cross-region backup.

* `status` - Indicates the status of the cross-region backup.

* `instance_id` - Indicates the ID of the RDS Instance.

* `begin_time` - Indicates the backup start time in the **yyyy-mm-ddThh:mm:ssZ** format

* `end_time` - Indicates the backup end time in the **yyyy-mm-ddThh:mm:ssZ"** format.

* `size` - Indicates the backup size in KB.

* `databases` - Indicates the database to be backed up.

  The [databases](#backups_databases_struct) structure is documented below.

* `associated_with_ddm` - Indicates whether a DDM instance has been associated.

* `datastore` - Indicates the database information

  The [datastore](#backups_datastore_struct) structure is documented below.

<a name="backups_databases_struct"></a>
The `databases` block supports:

* `name` - Indicates the database to be backed up for Microsoft SQL Server.

<a name="backups_datastore_struct"></a>
The `datastore` block supports:

* `type` - Indicates the database engine.
  Its value can be any of the following and is case-insensitive: **MySQL**, **PostgreSQL**, **SQLServer** and **MariaDB**.

* `version` - Indicates the database engine version.
