---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_backup"
description: |-
  Manages a GeminiDB manual backup resource within HuaweiCloud.
---

# huaweicloud_geminidb_backup

Manages a GeminiDB manual backup resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_geminidb_backup" "test" {
  instance_id = var.instance_id
  name        = "test-backup"
  description = "test backup description"
}
```

### With Database Tables (GeminiDB Cassandra Only)

```hcl
variable "instance_id" {}

resource "huaweicloud_geminidb_backup" "test" {
  instance_id = var.instance_id
  name        = "test-backup"
  description = "test backup with database tables"

  database_tables {
    database_name = "test_db"
    table_names   = ["table1", "table2"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GeminiDB instance.

* `name` - (Required, String, NonUpdatable) Specifies the name of the manual backup.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the manual backup.

* `database_tables` - (Optional, List, NonUpdatable) Specifies the database table information for the backup.
  The [database_tables](#geminidb_backup_database_tables) structure is documented below.

<a name="geminidb_backup_database_tables"></a>
The `database_tables` block supports:

* `database_name` - (Required, String, NonUpdatable) Specifies the database name.

* `table_names` - (Optional, List, NonUpdatable) Specifies the list of table names.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `backup_id` - The backup ID.

* `status` - The backup status.
  The valid values are as follows:
  + **BUILDING**: Backup in progress.
  + **COMPLETED**: Backup completed.
  + **FAILED**: Backup failed.

* `type` - The backup type.
  The valid values are as follows:
  + **Auto**: Automatic full backup.
  + **Manual**: Manual full backup.
  + **Incremental**: Incremental backup.
  + **Differential**: Differential backup.

* `size` - The backup size in KB.

* `begin_time` - The backup start time in UTC format "yyyy-mm-dd hh:mm:ss".

* `end_time` - The backup end time in UTC format "yyyy-mm-dd hh:mm:ss".

* `datastore` - The database information.
  The [datastore](#geminidb_backup_datastore) structure is documented below.

<a name="geminidb_backup_datastore"></a>
The `datastore` block supports:

* `type` - The database engine type.

* `version` - The database engine version.

## Import

The GeminiDB backup can be imported using the `backup_id`, e.g.

```
$ terraform import huaweicloud_geminidb_backup.test <backup_id>
```

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 5 minutes.
