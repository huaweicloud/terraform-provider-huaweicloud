---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_slow_log_details"
description: |-
  Use this data source to get the list of DAS slow log details.
---

# huaweicloud_das_slow_log_details

Use this data source to get the list of DAS slow log details.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

data "huaweicloud_das_slow_log_details" "test" {
  instance_id = var.instance_id
  start_time  = "2026-06-01T00:00:00+08:00"
  end_time    = "2026-06-02T00:00:00+08:00"
}
```

### Filter with execute time

```hcl
variable "instance_id" {}

data "huaweicloud_das_slow_log_details" "test" {
  instance_id      = var.instance_id
  start_time       = "2026-06-01T00:00:00+08:00"
  end_time         = "2026-06-02T00:00:00+08:00"
  execute_time_min = 10
  execute_time_max = 100
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the slow log details are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `start_time` - (Required, String) Specifies the start time of the query range, in RFC3339 format.

* `end_time` - (Required, String) Specifies the end time of the query range, in RFC3339 format.

-> 1.The earliest `start_time` is at most two days earlier than the current time.  
2.The latest `end_time` is at most one day later than the current time.  
3.The `end_time` must be greater than the `start_time`.

* `node_ids` - (Optional, List) Specifies the list of node IDs.

* `db_name` - (Optional, String) Specifies the database name.

* `sort_field` - (Optional, String) Specifies the field to sort by.  
  The valid values are as follows:
  + **occurrenceTime**. Execution start time.
  + **executeTime**. Execution loss time.
  + **lockWaitTime**. Lock wait time.
  + **rowsExamined**. The scanned row.
  + **rowsSent**. The returned row.

* `sort_asc` - (Optional, Bool) Specifies the sort order.  
  The valid values are as follows:
  + **true**. Ascending.
  + **false**. Descending.

* `client_ip_address` - (Optional, String) Specifies the client IP address.

* `user_name` - (Optional, String) Specifies the user name.

* `killed` - (Optional, String) Specifies the execution status.

* `execute_time_min` - (Optional, Int) Specifies the minimum execution time, in milliseconds.

* `execute_time_max` - (Optional, Int) Specifies the maximum execution time, in milliseconds.

* `rows_max_examined` - (Optional, Int) Specifies the maximum number of rows examined.

* `rows_min_examined` - (Optional, Int) Specifies the minimum number of rows examined.

* `fuzzy_sql` - (Optional, String) Specifies the fuzzy SQL pattern.

* `operation` - (Optional, String) Specifies the SQL operation types, separated by commas.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `details` - The list of slow log details that matched the filter parameters.  
  The [details](#slow_log_details_attr) structure is documented below.

<a name="slow_log_details_attr"></a>
The details block supports:

* `occurrence_time` - The occurrence time, in RFC3339 format.

* `sql_template_id` - The SQL template ID.

* `original_sql` - The original SQL statement.

* `db_name` - The database name.

* `client_ip_address` - The client IP address.

* `user_name` - The user name.

* `execute_time` - The execution time, in seconds.

* `lock_wait_time` - The lock wait time, in seconds.

* `rows_examined` - The number of rows examined.

* `rows_sent` - The number of rows sent.

* `tunable` - Whether the SQL can be tuned.

* `end_time` - The end time, in RFC3339 format (for SQL Server).

* `app_name` - The application name (for SQL Server).

* `rows_affected` - The number of rows affected (for SQL Server).

* `cpu_time` - The CPU time, in milliseconds (for SQL Server).

* `logical_reads` - The number of logical reads (for SQL Server).

* `physical_reads` - The number of physical reads (for SQL Server).

* `writes` - The number of writes (for SQL Server).

* `sql_type` - The SQL type.

* `collection` - The collection name (for MongoDB).

* `key_examined` - The number of keys examined (for MongoDB).

* `node_id` - The node ID (for MongoDB).

* `node_name` - The node name (for MongoDB).

* `killed` - The execution status.
