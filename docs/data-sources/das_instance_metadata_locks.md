---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_instance_metadata_locks"
description: |-
  Use this data source to get the list of DAS instance metadata locks.
---

# huaweicloud_das_instance_metadata_locks

Use this data source to get the list of DAS instance metadata locks.

-> This data source only supports to query metadata locks of **MySQL** instances.

## Example Usage

### Basic usage

```hcl
variable "instance_id" {}
variable "db_user_id" {}

data "huaweicloud_das_instance_metadata_locks" "test" {
  instance_id = var.instance_id
  db_user_id  = var.db_user_id
}
```

### Query metadata locks with database and table name

```hcl
variable "instance_id" {}
variable "db_user_id" {}
variable "database_name" {}
variable "table_name" {}

data "huaweicloud_das_instance_metadata_locks" "test" {
  instance_id = var.instance_id
  db_user_id  = var.db_user_id
  database    = var.database_name
  table       = var.table_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the instance metadata locks are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `db_user_id` - (Required, String) Specifies the database user ID.

* `thread_id` - (Optional, String) Specifies the session ID for filtering.

* `database` - (Optional, String) Specifies the database name for filtering.

* `table` - (Optional, String) Specifies the table name for filtering.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `metadata_locks` - The list of metadata locks that matched filter parameters.  
  The [metadata_locks](#instance_metadata_locks_metadata_locks) structure is documented below.

<a name="instance_metadata_locks_metadata_locks"></a>
The `metadata_locks` block supports:

* `thread_id` - The session ID of the thread holding or waiting for the metadata lock.

* `lock_status` - The lock status.
  + **PENDING**: Waiting for lock.
  + **GRANTED**: Holding lock.

* `lock_mode` - The lock mode.
  + **MDL_SHARED**
  + **MDL_EXCLUSIVE**
  + **MDL_SHARED_READ**
  + **MDL_SHARED_WRITE**

* `lock_type` - The lock type.
  + **Table metadata lock**
  + **Schema metadata lock**
  + **Tablespace lock**
  + **Global read lock**

* `lock_duration` - The lock duration.
  + **MDL_STATEMENT**: Statement level.
  + **MDL_TRANSACTION**: Transaction level.
  + **MDL_EXPLICIT**: Global level.

* `table_schema` - The database where the lock is held.

* `table_name` - The name of the table on which the lock is held.

* `user` - The database user associated with the session.

* `time` - The session duration.

* `host` - The host from which the session is connected.

* `database` - The database name of the session.

* `command` - The command being executed by the session.

* `state` - The state of the session.

* `sql` - The SQL statement being executed by the session.

* `trx_exec_time` - The execution time of the current transaction.

* `block_process` - The list of processes that are blocking the current lock.  
  The [block_process](#instance_metadata_locks_processes) structure is documented below.

* `wait_process` - The list of processes that are waiting for the current lock.  
  The [wait_process](#instance_metadata_locks_processes) structure is documented below.

<a name="instance_metadata_locks_processes"></a>
The `block_process` and `wait_process` block supports:

* `id` - The session ID of the process.

* `user` - The database user associated with the process.

* `host` - The connecting host of the process.

* `database` - The database name of the process.

* `command` - The command being executed by the process.

* `time` - The session duration of the process.

* `state` - The state of the process.

* `sql` - The SQL statement being executed by the process.

* `trx_executed_time` - The transaction duration of the process.
