---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_slow_log_statistics"
description: |-
  Use this data source to get the statistics of RDS slow logs within HuaweiCloud.
---

# huaweicloud_rds_slow_log_statistics

Use this data source to get the statistics of RDS slow logs within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_rds_slow_log_statistics" "test" {
  instance_id = var.instance_id
  start_time  = var.start_time
  end_time    = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `start_time` - (Required, String) Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - (Required, String) Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format.
  Only slow logs within one month before the current time can be queried.

* `type` - (Optional, String) Specifies the type of SQL statements.
  Value options: **INSERT**, **UPDATE**, **SELECT**, **DELETE**, **CREATE** and **ALL**, defaults to **ALL**.

* `database` - (Optional, String) Specifies the name of the database.
  The database name does not support searching with special characters: **<**, **>** and **&**.

* `sort` - (Optional, String) Specifies the sort field. If the value is **executeTime**, the results are sorted
  by the average execution time. Other values or an empty value sort the results by the number of executions.

* `order` - (Optional, String) Specifies the sort order. Value options: **desc** (descending) and
  **asc** (ascending), defaults to **desc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `slow_log_list` - The statistics list of the slow logs.

  The [slow_log_list](#slow_log_list_struct) structure is documented below.

<a name="slow_log_list_struct"></a>
The `slow_log_list` block supports:

* `count` - The number of executions.

* `time` - The average execution time.

* `lock_time` - The average lock wait time, valid only for MySQL.

* `rows_sent` - The average number of result rows, valid only for MySQL.

* `rows_examined` - The average number of scanned rows, valid only for MySQL.

* `database` - The database name.

* `users` - The account name.

* `query_sample` - The SQL execution syntax.

* `client_ip` - The IP address of the client.

* `type` - The type of SQL statements.
