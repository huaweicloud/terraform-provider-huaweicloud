---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_replication_db_parameters"
description: |-
  Use this data source to query the database parameters configuration of StarRocks data replication within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_replication_db_parameters

Use this data source to query the database parameters configuration of StarRocks data replication within HuaweiCloud.

## Example Usage

```hcl
variable "htap_instance_id" {}

data "huaweicloud_taurusdb_htap_starrocks_replication_db_parameters" "test" {
  instance_id = var.htap_instance_id
}
```

### Query with filters

```hcl
variable "htap_instance_id" {}

data "huaweicloud_taurusdb_htap_starrocks_replication_db_parameters" "test" {
  instance_id       = var.htap_instance_id
  add_task_scenario = "add_sub_task"
  main_task_name    = "test_task"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the replication database parameters.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the StarRocks instance ID.

* `add_task_scenario` - (Optional, String) Specifies the scenario for adding a sub-task,
  used to distinguish whether the database parameter supports modification.

* `main_task_name` - (Optional, String) Specifies the main task name corresponding to the sub-task.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `db_parameters` - The list of database parameter configurations.
  The [db_parameters](#htap_starrocks_replication_db_parameters_attr) structure is documented below.

<a name="htap_starrocks_replication_db_parameters_attr"></a>
The `db_parameters` block supports:

* `param_name` - The parameter name.

* `data_type` - The parameter data type. The valid values are **Integer** and **String**.

* `default_value` - The default value of the parameter.

* `value_range` - The value range of the parameter.

* `description` - The description of the parameter.

* `is_modifiable` - Whether the parameter is modifiable in the current scenario.
