---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_replication"
description: |-
  Manages a TaurusDB HTAP StarRocks data replication task resource within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_replication

Manages a TaurusDB HTAP StarRocks data replication task resource within HuaweiCloud.

## Example Usage

### Create replication task with instance all databases

```hcl
variable "htap_instance_id" {}
variable "taurusdb_instance_id" {}

resource "huaweicloud_taurusdb_htap_starrocks_replication" "test" {
  instance_id            = var.htap_instance_id
  source_instance_id     = var.taurusdb_instance_id
  task_name              = "instance_all_dbs_repl"
  is_instance_level_sync = "true"
  database_repl_scope    = "all"
  source_database_name   = "ALL"
  target_database_name   = "ALL"
  enable_sync            = "true"
  sync_action            = "pause"

  db_configs {
    param_name = "binlog_expire_logs_seconds"
    value      = "0"
  }

  db_configs {
    param_name = "max_full_sync_task_threads_num"
    value      = "4"
  }

  table_repl_config {
    repl_type  = "include_tables"
    repl_scope = "all"
  }
}
```

### Create replication task with specific database and tables

```hcl
variable "htap_instance_id" {}
variable "taurusdb_instance_id" {}
variable "taurusdb_slave_node_id" {}

resource "huaweicloud_taurusdb_htap_starrocks_replication" "test" {
  task_name              = "part_db_tables_repl"
  instance_id            = var.htap_instance_id
  source_instance_id     = var.taurusdb_instance_id
  source_node_id         = var.taurusdb_slave_node_id
  is_instance_level_sync = "false"
  database_repl_scope    = "part"
  source_database_name   = "__taurus_sys__"
  target_database_name   = "__taurus_sys__"

  db_configs {
    param_name = "binlog_expire_logs_seconds"
    value      = "0"
  }

  db_configs {
    param_name = "max_full_sync_task_threads_num"
    value      = "4"
  }

  table_repl_config {
    repl_type  = "include_tables"
    repl_scope = "part"
    tables     = ["tenant", "tenant_db"]
  }

  tables_configs {
    table_name   = "tenant"
    table_config = "order by tenant_name;key columns tenant_name"
  }

  tables_configs {
    table_name   = "tenant_db"
    table_config = "order by tenant_name,db_name"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the HTAP StarRocks replication task.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NoneUpdatable) Specifies the HTAP StarRocks instance ID.

* `task_name` - (Required, String, NoneUpdatable) Specifies the name of the data synchronization task.
  The value can contain 3 to 128 characters. Only uppercase letters, lowercase letters, digits, and
  underscores (_) are allowed.

* `source_instance_id` - (Required, String, NoneUpdatable) Specifies the source TaurusDB instance ID.

* `source_node_id` - (Optional, String) Specifies the TaurusDB read replica node ID.

* `is_instance_level_sync` - (Optional, String) Specifies whether instance-level synchronization is supported.
  Valid values are **true** and **false**. Defaults to **false**.

* `database_repl_scope` - (Optional, String) Specifies the database synchronization scope.
  Valid values are **all** (all databases are synchronized) and **part** (some databases are synchronized).
  Defaults to **part**.

* `source_database_name` - (Required, String) Specifies the name of the source database of the source TaurusDB instance.
  The value can contain 3 to 1,024 characters. Only uppercase letters, lowercase letters, digits, and
  underscores (_) are allowed.

* `target_database_name` - (Required, String) Specifies the name of the destination database.
  The value can contain 3 to 128 characters. Only uppercase letters, lowercase letters, digits, and
  underscores (_) are allowed.

* `db_configs` - (Required, List) Specifies the database configurations.
  The [db_configs](#replication_db_configs_attr) structure is documented below.

* `table_repl_config` - (Required, List) Specifies the table synchronization configurations.
  The [table_repl_config](#replication_table_repl_config_attr) structure is documented below.

* `tables_configs` - (Optional, List) Specifies the table configurations.
  The [tables_configs](#replication_tables_configs_attr) structure is documented below.

* `enable_sync` - (Optional, String) Specifies whether to enable the created data synchronization task.
  The valid values are as follows:
  + **true**: Start the data synchronization task.
  + **false**: Do not start the data synchronization task.

* `sync_action` - (Optional, String) Specifies the operation to be performed on the data synchronization task.
  The valid values are as follows:
  + **pause**: Pause the started data synchronization task.
  + **resume**: Resume the paused data synchronization task.

<a name="replication_db_configs_attr"></a>
The `db_configs` block supports:

* `param_name` - (Required, String) Specifies the parameter name.

* `value` - (Required, String) Specifies the parameter value.

<a name="replication_table_repl_config_attr"></a>
The `table_repl_config` block supports:

* `repl_type` - (Required, String) Specifies the table synchronization type.
  Valid values are **include_tables** (whitelist) and **exclude_tables** (blacklist).

* `repl_scope` - (Required, String) Specifies the table synchronization scope.
  Valid values are **all** (all tables) and **part** (some tables).

* `tables` - (Optional, List) Specifies the whitelisted or blacklisted tables.

<a name="replication_tables_configs_attr"></a>
The `tables_configs` block supports:

* `table_name` - (Required, String) Specifies the table name.

* `table_config` - (Required, String) Specifies the table configuration value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<instance_id>/<task_name>`.

* `status` - The current status of the data synchronization task.
  The valid values are as follows:
  + **Yes**: The task is normal.
  + **No**: The task is abnormal.

* `stage` - The stage of the data synchronization.
  The valid values are as follows:
  + **Wait**: Waiting for synchronization.
  + **Incremental**: Incremental synchronization.
  + **Full**: Full synchronization.
  + **Cancelled**: Synchronization cancelled.
  + **Paused**: Synchronization paused.

* `percentage` - The progress percentage. Valid range is `0` to `100`.

* `is_need_repair` - Whether the task need to be repaired.

* `is_main_task` - Whether the task is the main task.

* `database_info` - The TaurusDB database configuration information.
  The [database_info](#replication_database_info_attr) structure is documented below.

* `table_infos` - The table configuration check results.
  The [table_infos](#replication_table_infos_attr) structure is documented below.

* `new_table_repl_config` - The updated table synchronization configurations.
  The [new_table_repl_config](#replication_new_table_repl_config_attr) structure is documented below.

* `is_instance_level_sync` - Whether instance-level synchronization is supported.

* `is_support_reg_exp` - Whether wildcards are supported.

* `is_need_repair` - Whether the replication task needs repair.

* `is_main_task` - Whether this is the main replication task.

* `is_tables_change` - Whether there is any change to the synchronization scope.

* `error_msg` - The error message if the replication task failed.

* `last_error_of_alter_table` - Exception about the latest ALTER TABLE operation.

<a name="replication_database_info_attr"></a>
The `database_info` block supports:

* `database_name` - The database name.

* `db_config_check_results` - The database configuration check results.
  The [db_config_check_results](#replication_db_config_check_results_attr) structure is documented below.

<a name="replication_db_config_check_results_attr"></a>
The `db_config_check_results` block supports:

* `param_name` - The parameter name.

* `value` - The parameter value.

* `check_result` - The check result. Valid values are **success** and **fail**.

<a name="replication_table_infos_attr"></a>
The `table_infos` block supports:

* `table_name` - The table name.

* `table_config` - The table configuration item.

* `check_result` - The check result. Valid values are **success** and **failed**.

<a name="replication_new_table_repl_config_attr"></a>
The `new_table_repl_config` block supports:

* `repl_type` - The table synchronization type.

* `repl_scope` - The table synchronization scope.

* `tables` - The whitelisted or blacklisted tables.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The StarRocks replication task can be imported using the `instance_id` and `task_name`, separated by a slash, e.g.

```
$ terraform import huaweicloud_taurusdb_htap_starrocks_replication.test <instance_id>/<task_name>
```

Note that the imported state may not be identical to your resource definition, due to the attribute missing from the
API response. The missing attributes are: `db_configs`, `tables_configs`, `enable_sync`, `sync_action`. It is generally
recommended running `terraform plan` after importing a HTAP StarRocks replication task. You can then decide if changes
should be applied to the HTAP StarRocks replication task, or the resource definition should be updated to align with
the HTAP StarRocks replication task. Also you can ignore changes as below.

```hcl
resource "huaweicloud_taurusdb_htap_starrocks_replication" "test" {
  ...

  lifecycle {
    ignore_changes = [
      db_configs, tables_configs, enable_sync, sync_action
    ]
  }
}
```
