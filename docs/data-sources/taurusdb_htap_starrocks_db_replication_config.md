---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_db_replication_config"
description: |-
  Use this data source to query StarRocks data synchronization configurations by destination database within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_db_replication_config

Use this data source to query StarRocks data synchronization configurations by destination database within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "database" {}

data "huaweicloud_taurusdb_htap_starrocks_db_replication_config" "test" {
  instance_id = var.instance_id
  database    = var.database
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the HTAP StarRocks instance ID.

* `database` - (Required, String) Specifies the name of the destination database.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `source_instance_id` - The TaurusDB instance ID.

* `source_node_id` - The TaurusDB node ID.

* `is_instance_level_sync` - Whether instance-level synchronization is supported.
  Valid values are **true** and **false**.

* `database_repl_scope` - The database synchronization scope.
  Valid values are **all** (all databases are synchronized) and **part** (some databases are synchronized).

* `target_database_name` - The name of the destination database.

* `is_support_reg_exp` - Whether wildcards are supported.

* `is_tables_change` - Whether there is any change to the synchronization scope.

* `error_msg` - The error message if the replication task failed.

* `last_error_of_alter_table` - Exception about the latest ALTER TABLE operation.

* `database_info` - The TaurusDB database configuration.
  The [database_info](#database_info_struct) structure is documented below.

* `table_infos` - The table configurations.
  The [table_infos](#table_infos_struct) structure is documented below.

* `table_repl_config` - The table synchronization configurations.
  The [table_repl_config](#table_repl_config_struct) structure is documented below.

* `new_table_repl_config` - The updated table synchronization configurations.
  The [new_table_repl_config](#table_repl_config_struct) structure is documented below.

<a name="database_info_struct"></a>
The `database_info` block supports:

* `database_name` - The database name.

* `db_config_check_results` - The database configuration check results.
  The [db_config_check_results](#db_config_check_results_struct) structure is documented below.

<a name="db_config_check_results_struct"></a>
The `db_config_check_results` block supports:

* `param_name` - The parameter name.

* `value` - The parameter value.

* `check_result` - The check result. Valid values are **success** and **fail**.

<a name="table_infos_struct"></a>
The `table_infos` block supports:

* `table_name` - The table name.

* `table_config` - The table configuration item.

* `check_result` - The check result. Valid values are **success** and **failed**.

<a name="table_repl_config_struct"></a>
The `table_repl_config` and `new_table_repl_config` blocks support:

* `repl_type` - The table synchronization type.
  Valid values are **include_tables** (whitelist) and **exclude_tables** (blacklist).

* `repl_scope` - The table synchronization scope. Valid values are **all** (all tables) and **part** (some tables).

* `tables` - The whitelisted or blacklisted tables.
