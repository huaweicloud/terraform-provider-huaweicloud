---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_sql_throttling_task"
description: |-
  Manages a GaussDB OpenGauss SQL throttling task resource within HuaweiCloud.
---

# huaweicloud_gaussdb_opengauss_sql_throttling_task

Manages a GaussDB OpenGauss SQL throttling task resource within HuaweiCloud.

## Example Usage

### Create by SQL_ID

```hcl
variable "instance_id"{}
variable "node_id"{}
variable "sql_id"{}

resource "huaweicloud_gaussdb_opengauss_sql_throttling_task" "source" {
  instance_id      = var.instance_id
  task_scope       = "SQL"
  limit_type       = "SQL_ID"
  limit_type_value = var.sql_id
  task_name        = "test_task_name"
  parallel_size    = 4
  start_time       = "2025-02-12T11:47:20+0800"
  end_time         = "2025-02-13T11:47:20+0800"
  sql_model        = "select name, setting from pg_settings where name in (1...n)"

  node_infos {
    node_id = var.node_id
    sql_id  = var.sql_id
  }
}
```

### Create by SQL_TYPE

```hcl
variable "instance_id"{}

resource "huaweicloud_gaussdb_opengauss_sql_throttling_task" "test" {
  instance_id      = var.instance_id
  task_scope       = "SQL"
  limit_type       = "SQL_TYPE"
  limit_type_value = "update"
  task_name        = "test_task_name"
  parallel_size    = 4
  start_time       = "2025-02-12T11:47:20+0800"
  end_time         = "2025-02-13T11:47:20+0800"
  key_words        = "aaa,bbb,ccc"
  databases        = "test_db_111"
}
```

### Create by SESSION

```hcl
variable "instance_id"{}

resource "huaweicloud_gaussdb_opengauss_sql_throttling_task" "test" {
  instance_id        = var.instance_id
  task_scope         = "SESSION"
  limit_type         = "SESSION_ACTIVE_MAX_COUNT"
  limit_type_value   = "CPU_OR_MEMORY"
  task_name          = "test_task_name"
  parallel_size      = 4
  cpu_utilization    = 20
  memory_utilization = 40
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB OpenGauss instance.

  Changing this parameter will create a new resource.

* `task_scope` - (Required, String, ForceNew) Specifies the task scope. Currently, **SQL** and **SESSION** are supported.

  Changing this parameter will create a new resource.

* `limit_type` - (Required, String, ForceNew) Specifies the throttling type.
  + When `task_scope` is set to **SQL**, the value can be **SQL_ID** or **SQL_TYPE**.
  + When `task_scope` is set to **SESSION**, the value can be **SESSION_ACTIVE_MAX_COUNT**.

  Changing this parameter will create a new resource.

* `limit_type_value` - (Required, String, ForceNew) Specifies the throttling type value.
  + When `limit_type` is set to **SQL_ID**, the value of this parameter is the SQL ID of the selected template.
  + When `limit_type` is set to **SQL_TYPE**, the value of this parameter can be **select**, **update**, **insert**,
    **delete**, or **merge**.
  + When `limit_type` is set to **SESSION_ACTIVE_MAX_COUNT**, the value of this parameter can only be **CPU_OR_MEMORY**.

  Changing this parameter will create a new resource.

* `task_name` - (Required, String, ForceNew) Specifies the name of the SQL throttling task. The value can contain up to
  **100** characters. Only uppercase letters, lowercase letters, underscores (_), digits, and dollar signs ($) are allowed.

* `parallel_size` - (Required, Int) Specifies the maximum concurrency. The value can be **0** or a positive integer.
  Value range: **0** to **2147483647**.

* `start_time` - (Optional, String) Specifies the task start time. It is mandatory when `task_scope` is set to **SQL**.
  Value range: two minutes later than or equal to the current time (UTC time). The format must be **yyyy-mm-ddThh:mm:ss+0000**.

* `end_time` - (Optional, String) Specifies the task end time. It is mandatory when `task_scope` is set to **SQL**.
  Value range: later than the task start time. The format must be **yyyy-mm-ddThh:mm:ss+0000**.

* `key_words` - (Optional, String) Specifies the keyword. It is mandatory when `limit_type` is set to **SQL_TYPE**. You
  can enter **2** to **100** keywords and separate multiple keywords by commas (,). Each keyword can contain **2** to
  **64** characters and cannot start and end with a space. The specifical characters ("\{}) and null are not allowed.

* `sql_model` - (Optional, String, ForceNew) Specifies the SQL template. It is mandatory when `limit_type` is set to
  **SQL_ID**.

  Changing this parameter will create a new resource.

* `cpu_utilization` - (Optional, Int) Specifies the CPU usage threshold. The value is an integer ranging from **0** to
  **100**. It is mandatory when `limit_type` is set to **SESSION_ACTIVE_MAX_COUNT**. This parameter and `memory_utilization`
  cannot be both set to **0**. If you only need one of them for throttling, set the other threshold to **0**.

* `memory_utilization` - (Optional, Int) Specifies the Memory usage threshold. The value is an integer ranging from **0**
  to **100**. It is mandatory when `limit_type` is set to **SESSION_ACTIVE_MAX_COUNT**. This parameter and
  `cpu_utilization` cannot be both set to **0**. If you only need one of them for throttling, set the other threshold to
  **0**.

* `databases` - (Optional, String) Specifies the databases of the instance. Databases are separated by commas (,). It is
  mandatory when `limit_type` is set to **SQL_TYPE**.

* `node_infos` - (Optional, List, ForceNew) Specifies the CN information. It is mandatory when `limit_type` is set to
  **SQL_ID**. The [node_infos](#node_infos_struct) structure is documented below.

  Changing this parameter will create a new resource.

<a name="node_infos_struct"></a>
The `node_infos` block supports:

* `node_id` - (Required, String, ForceNew) Specifies the node ID.

  Changing this parameter will create a new resource.

* `sql_id` - (Required, String, ForceNew) Specifies the ID of the SQL statement executed on the node. If `limit_type` is
  set to **SQL_ID**, the value of this parameter must be the same as that of `limit_type_value`.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - Indicates the creation time in the format of **yyyy-mm-ddThh:mm:ssZ**.

* `updated_at` - Indicates the update time in the format of **yyyy-mm-ddThh:mm:ssZ**.

* `creator` - Indicates the creator.

* `modifier` - Indicates the modifier.

* `rule_name` - Indicates the rule name.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 90 minutes.
* `update` - Default is 90 minutes.
* `delete` - Default is 90 minutes.

## Import

The GaussDB OpenGauss SQL throttling task can be imported using `instance_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_gaussdb_opengauss_sql_throttling_task.test <instance_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to the attribute missing from the
API response. The missing attribute is: `start_time` `end_time`. It is generally recommended running `terraform plan`
after importing a GaussDB OpenGauss SQL throttling task. You can then decide if changes should be applied to the GaussDB
OpenGauss SQL throttling task, or the resource definition should be updated to align with the GaussDB OpenGauss SQL
throttling task. Also you can ignore changes as below.

```hcl
resource "huaweicloud_gaussdb_opengauss_sql_throttling_task" "test" {
  ...

  lifecycle {
    ignore_changes = [
      start_time, end_time,
    ]
  }
}
```
