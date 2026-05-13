---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_backups"
description: |-
  Use this data source to get a list of GeminiDB backups.
---

# huaweicloud_geminidb_backups

Use this data source to get a list of GeminiDB backups.

## Example Usage

### List All Backups

```hcl
data "huaweicloud_geminidb_backups" "test" {}
```

### Filter by Instance ID

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_backups" "test" {
  instance_id = var.instance_id
}
```

### Filter by Datastore Type

```hcl
data "huaweicloud_geminidb_backups" "test" {
  datastore_type = "redis"
}
```

### Filter by Backup Type

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_backups" "test" {
  instance_id = var.instance_id
  backup_type = "Auto"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the backups.
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the ID of the GeminiDB instance.

* `datastore_type` - (Optional, String) Specifies the database type.
  The valid values are as follows:
  + **cassandra**: GeminiDB Cassandra.
  + **redis**: GeminiDB Redis.
  + **mongodb**: GeminiDB Mongo.
  + **influxdb**: GeminiDB Influx.

* `backup_id` - (Optional, String) Specifies the backup ID.

* `backup_type` - (Optional, String) Specifies the backup type.
  The valid values are as follows:
  + **Auto**: Automatic full backup.
  + **Manual**: Manual full backup.

* `type` - (Optional, String) Specifies the backup strategy type.
  The valid values are as follows:
  + **Instance**: Instance-level backup.
  + **DatabaseTable**: Database-table-level backup.

* `begin_time` - (Optional, String) Specifies the start time for the query.
  The format is `yyyy-mm-ddThh:mm:ssZ`.
  When `end_time` is specified, `begin_time` is required.

* `end_time` - (Optional, String) Specifies the end time for the query.
  The format is `yyyy-mm-ddThh:mm:ssZ`.
  When `begin_time` is specified, `end_time` is required.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `backups` - The list of backups.
  The [backups](#geminidb_backups_backup) structure is documented below.

<a name="geminidb_backups_backup"></a>
The `backups` block supports:

* `id` - The backup ID.

* `name` - The backup name.

* `instance_id` - The ID of the instance to which the backup belongs.

* `instance_name` - The name of the instance to which the backup belongs.

* `datastore` - The database version information.
  The [datastore](#geminidb_backups_datastore) structure is documented below.

* `type` - The backup type.
  + **Auto**: Automatic full backup.
  + **Manual**: Manual full backup.

* `size` - The backup size in KB.

* `status` - The backup status.
  + **BUILDING**: Backup in progress.
  + **COMPLETED**: Backup completed.
  + **FAILED**: Backup failed.

* `begin_time` - The backup start time.

* `end_time` - The backup end time.

* `description` - The backup description.

* `database_tables` - The database table information in the backup.
  The [database_tables](#geminidb_backups_database_tables) structure is documented below.
  This field is empty for instance-level queries and can be ignored.
  For database-table-level queries, this field is non-empty if there are database-table-level backups.

<a name="geminidb_backups_datastore"></a>
The `datastore` block supports:

* `type` - The database type.

* `version` - The database version.

<a name="geminidb_backups_database_tables"></a>
The `database_tables` block supports:

* `database_name` - The database name.

* `table_names` - The list of table names.
  When `table_names` is empty, it indicates a database-level query.
  When `table_names` is non-empty, it indicates a table-level query.
