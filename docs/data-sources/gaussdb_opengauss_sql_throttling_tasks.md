---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_sql_throttling_tasks"
description: |-
  Use this data source to get the SQL throttling tasks based on search criteria.
---

# huaweicloud_gaussdb_opengauss_sql_throttling_tasks

Use this data source to get the SQL throttling tasks based on search criteria.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_opengauss_sql_throttling_tasks" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB OpenGauss instance.

* `task_scope` - (Optional, String) Specifies the throttling task scope.
  Currently, **SQL** and **SESSION** are supported.

* `limit_type` - (Optional, String) Specifies the throttling type.
  The value can be **SQL_ID**, **SQL_TYPE** or **SESSION_ACTIVE_MAX_COUNT**.

* `limit_type_value` - (Optional, String) Specifies the throttling type value. Fuzzy match is supported.

* `task_name` - (Optional, String) Specifies the throttling task name. Fuzzy match is supported.

* `sql_model` - (Optional, String) Specifies the SQL template. Fuzzy match is supported.

* `rule_name` - (Optional, String) Specifies the rule name.

* `start_time` - (Optional, String) Specifies the start time of the throttling task in the format of **yyy-mm-ddThh:mm:ss+0000**.

* `end_time` - (Optional, String) Specifies the end time of the throttling task in the format of **yyy-mm-ddThh:mm:ss+0000**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `limit_task_list` - Indicates the list of throttling tasks.

  The [limit_task_list](#limit_task_list_struct) structure is documented below.

<a name="limit_task_list_struct"></a>
The `limit_task_list` block supports:

* `instance_id` - Indicates the instance ID.

* `task_id` - Indicates the throttling task ID.

* `task_name` - Indicates the throttling task name.

* `task_scope` - Indicates the throttling task scope.

* `limit_type` - Indicates the throttling task type.

* `limit_type_value` - Indicates the throttling task type value.

* `sql_model` - Indicates the SQL template.
  This parameter is returned only when the `limit_type` is **SQL_ID**.

* `key_words` - Indicates the keyword.
  This parameter is returned only when the `limit_type` is **SQL_TYPE**.

* `status` - Indicates the throttling task status.
  The value can be **CREATING**, **UPDATING**, **DELETING**, **WAIT_EXCUTE**, **EXCUTING**, **TIME_OVER**, **DELETED**,
  **CREATE_FAILED**, **UPDATE_FAILED**, **DELETE_FAILED**, **EXCEPTION** or **NODE_SHUT_DOWN**.

* `rule_name` - Indicates the rule name.

* `parallel_size` - Indicates the maximum concurrency.

* `start_time` - Indicates the start time of the throttling task in the format of **yyyy-mm-ddThh:mm:ssZ**.

* `end_time` - Indicates the end time of the throttling task in the format of **yyyy-mm-ddThh:mm:ssZ**.

* `cpu_utilization` - Indicates the CPU usage.
  This parameter is returned only when the `limit_type` is **SESSION_ACTIVE_MAX_COUNT**.

* `memory_utilization` - Indicates the memory usage.
  This parameter is returned only when the `limit_type` is **SESSION_ACTIVE_MAX_COUNT**.

* `created_at` - Indicates the creation time in the format of **yyyy-mm-ddThh:mm:ssZ**.

* `updated_at` - Indicates the update time in the format of **yyyy-mm-ddThh:mm:ssZ**.

* `creator` - Indicates the creator.

* `modifier` - Indicates the modifier.

* `databases` - Indicates the databases of the instance.
  Databases are separated by commas (,).

* `node_infos` - Indicates the CN information.

  The [node_infos](#limit_task_list_node_infos_struct) structure is documented below.

<a name="limit_task_list_node_infos_struct"></a>
The `node_infos` block supports:

* `node_id` - Indicates the node ID.

* `sql_id` - Indicates the ID of the SQL statement executed on the node.
