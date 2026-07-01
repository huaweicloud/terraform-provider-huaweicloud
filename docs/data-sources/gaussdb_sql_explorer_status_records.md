---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_sql_explorer_status_records"
description: |-
  Use this data source to query full SQL switch configurations of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_sql_explorer_status_records

Use this data source to query full SQL switch configurations of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_sql_explorer_status_records" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `full_sql_switches` - The list of full SQL switch configuration records.
  The [full_sql_switches](#gaussdb_instance_full_sql_switches) structure is documented below.

* `allowed_sql_types` - The list of system preset SQL type rules supported by the instance.
  The [allowed_sql_types](#gaussdb_instance_allowed_sql_types) structure is documented below.

<a name="gaussdb_instance_full_sql_switches"></a>
The `full_sql_switches` block supports:

* `is_open` - Whether the full SQL collection task is enabled.

* `begin_time` - The start timestamp of the collection task.

* `end_time` - The end timestamp of the collection task, null if the task is running.

* `save_days` - The log retention days.

* `storage_mode` - The log storage mode, fixed to LTS.

* `is_exclude_sys_user` - Whether to exclude system user SQL records.

* `lts_config` - The LTS log storage configuration block.
  The [lts_config](#gaussdb_instance_lts_config) structure is documented below.

* `sql_type_range` - Custom SQL type filtering rules for this switch task.
  The [sql_type_range](#gaussdb_instance_allowed_sql_types) structure is documented below.

<a name="gaussdb_instance_lts_config"></a>
The `lts_config` block supports:

* `group_ttl_in_days` - The retention days of the LTS log group.

* `group_log_type` - The log type of the log group.

* `log_group_name` - The name of the LTS log group.

* `log_group_id` - The ID of the LTS log group.

* `log_stream_name` - The name of the LTS log stream.

* `log_stream_id` - The ID of the LTS log stream.

* `stream_log_type` - The log type of the log stream.

* `stream_ttl_in_days` - The retention days of the log stream.

* `stream_structure_config_id` - The ID of stream structure configuration.

* `stream_index_config_id` - The ID of stream index configuration.

<a name="gaussdb_instance_allowed_sql_types"></a>
The `sql_type_range` / `allowed_sql_types` block supports:

* `category` - The SQL rule category, valid values: `all`, `dml`, `ddl`, `dcl`, `tcl`, `dql`, `custom`.

* `prefixes` - The list of SQL statement prefix filters.

* `is_preset` - Whether this rule is a system built-in preset rule.
