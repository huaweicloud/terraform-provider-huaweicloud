---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_slow_log_templates"
description: |-
  Use this data source to get the list of DAS slow log templates.
---

# huaweicloud_das_slow_log_templates

Use this data source to get the list of DAS slow log templates.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

data "huaweicloud_das_slow_log_templates" "test" {
  instance_id = var.instance_id
  start_time  = "2026-06-08T09:00:00+08:00"
  end_time    = "2026-06-08T19:00:00+08:00"
}
```

### Filter with average execution time

```hcl
variable "instance_id" {}

data "huaweicloud_das_slow_log_templates" "test" {
  instance_id          = var.instance_id
  start_time           = "2026-06-08T09:00:00+08:00"
  end_time             = "2026-06-08T19:00:00+08:00"
  min_avg_execute_time = 10
  max_avg_execute_time = 100
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the slow log templates are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `start_time` - (Required, String) Specifies the start time of the query range, in RFC3339 format.

* `end_time` - (Required, String) Specifies the end time of the query range, in RFC3339 format.

-> 1.The maximum difference between `start_time` and `end_time` is `12` hours.  
2.The `end_time` must be greater than the `start_time`.  
3.The `start_time` and `end_time` only support **year**, **month**, **day**, and **hour**, but not **minute** and
**second**, such as **2026-06-08T09:00:00+08:00**.

* `template_id` - (Optional, String) Specifies the SQL template ID.

* `node_id` - (Optional, String) Specifies the node ID of the instance.

* `db_name` - (Optional, String) Specifies the database name.

* `min_avg_execute_time` - (Optional, Float) Specifies the minimum average execution time, in milliseconds.

* `max_avg_execute_time` - (Optional, Float) Specifies the maximum average execution time, in milliseconds.

* `operation` - (Optional, String) Specifies the SQL operation types, separated by commas.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - The list of slow log templates that matched filter parameters.  
  The [templates](#slow_log_templates_attr) structure is documented below.

<a name="slow_log_templates_attr"></a>
The templates block supports:

* `template_name` - The name of SQL template.

* `template_id` - The ID of the SQL template.

* `sql_sample` - The SQL sample.

* `db_names` - The database names.

* `execute_count` - The execution count.

* `avg_execute_time` - The average execution time, in milliseconds.

* `max_execute_time` - The maximum execution time, in milliseconds.

* `avg_lock_wait_time` - The average lock wait time, in milliseconds.

* `max_lock_wait_time` - The maximum lock wait time, in milliseconds.

* `avg_rows_examined` - The average number of rows examined.

* `max_rows_examined` - The maximum number of rows examined.

* `avg_rows_sent` - The average number of rows sent.

* `max_rows_sent` - The maximum number of rows sent.

* `tunable` - Whether the SQL can be tuned.

* `node_ids` - The node IDs.

* `avg_cpu_time` - The average CPU time, in milliseconds.

* `max_cpu_time` - The maximum CPU time, in milliseconds.

* `avg_rows_affected` - The average number of rows affected.

* `max_rows_affected` - The maximum number of rows affected.

* `avg_logical_reads` - The average number of logical reads.

* `max_logical_reads` - The maximum number of logical reads.

* `avg_physical_reads` - The average number of physical reads.

* `max_physical_reads` - The maximum number of physical reads.

* `avg_writes` - The average number of writes.

* `max_writes` - The maximum number of writes.

* `instance_id` - The instance ID.

* `total_execute_time_ratio` - The total execution time ratio.

* `execute_count_ratio` - The execution count ratio.
