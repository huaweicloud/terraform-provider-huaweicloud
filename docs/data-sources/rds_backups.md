---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_backups"
description: |-
  Use this data source to get the list of RDS backups.
---

# huaweicloud_rds_backups

Use this data source to get the list of RDS backups.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_backups" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Instance ID.

* `backup_id` - (Optional, String) Backup ID.

* `backup_type` - (Optional, String) Backup type.  
  The options are as follows:
  + **auto**: Automated full backup.
  + **manual**: Manual full backup.
  + **fragment**: Differential full backup.
  + **incremental**: Automated incremental backup.

* `begin_time` - (Optional, String) Start time in the "yyyy-mm-ddThh:mm:ssZ" format.

* `end_time` - (Optional, String) End time in the "yyyy-mm-ddThh:mm:ssZ" format.

* `name` - (Optional, String) Backup name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `backups` - Backup list. For details, see Data structure of the Backup field.
  The [backups](#Backup_Backup) structure is documented below.

<a name="Backup_Backup"></a>
The `backups` block supports:

* `id` - Backup ID.

* `instance_id` - RDS instance ID.

* `name` - Backup name.

* `type` - Backup type.  
  The options are as follows:
  + **auto**: Automated full backup.
  + **manual**: Manual full backup.
  + **fragment**: Differential full backup.
  + **incremental**: Automated incremental backup.

* `size` - Backup size in KB.

* `status` - Backup status.  
  The options are as follows:
  + **BUILDING**: Backup in progress.
  + **COMPLETED**: Backup completed.
  + **FAILED**: Backup failed.
  + **DELETING**: Backup being deleted.

* `begin_time` - Backup start time in the "yyyy-mm-ddThh:mm:ssZ" format.

* `end_time` - Backup end time in the "yyyy-mm-ddThh:mm:ssZ" format.

* `associated_with_ddm` - Whether a DDM instance has been associated.

* `datastore` - The database information.
  The [datastore](#Backup_BackupDatastore) structure is documented below.

* `databases` - Database been backed up.
  The [databases](#Backup_BackupDatabases) structure is documented below.

<a name="Backup_BackupDatastore"></a>
The `datastore` block supports:

* `type` - DB engine.  
The value can be **MySQL**, **PostgreSQL**, **SQLServer**, **MariaDB**.

* `version` - DB engine version.

<a name="Backup_BackupDatabases"></a>
The `databases` block supports:

* `name` - Database to be backed up for Microsoft SQL Server.
