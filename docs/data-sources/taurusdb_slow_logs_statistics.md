---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_slow_logs_statistics"
description: |-
  Use this data source to query the list of TaurusDB slow log statistics within HuaweiCloud.
---

# huaweicloud_taurusdb_slow_logs_statistics

Use this data source to query the list of TaurusDB slow log statistics within HuaweiCloud.

## Example Usage

### Query TaurusDB Slow Logs Statistics

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_taurusdb_slow_logs_statistics" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  start_time  = var.start_time
  end_time    = var.end_time
}
```

### Query TaurusDB Slow Logs Statistics Filtered by SQL Type and Database

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "start_time" {}
variable "end_time" {}
variable "type" {}
variable "database" {}

data "huaweicloud_taurusdb_slow_logs_statistics" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  start_time  = var.start_time
  end_time    = var.end_time
  type        = var.type
  database    = var.database
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the TaurusDB instance.

* `node_id` - (Required, String) Specifies the ID of the instance node.

* `start_time` - (Required, String) Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - (Required, String) Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `type` - (Optional, String) Specifies the statement type.
  The valid values are as follows:
  + **INSERT**
  + **UPDATE**
  + **SELECT**
  + **DELETE**
  + **CREATE**
  + **ALL**

* `database` - (Optional, String) Specifies the database name.

* `sort` - (Optional, String) Specifies the sorting field.
  If this parameter is set to **executeTime**, slow query logs are sorted by average execution duration.
  If this parameter is left empty or set to other values, the slow query logs are sorted by executions.

* `order` - (Optional, String) Specifies the sorting order.
  The valid values are as follows:
  + **desc**: The query results are displayed in descending order.
  + **asc**: The query results are displayed in ascending order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `slow_log_list` - Indicates the list of the slow log statistics.
  The [slow_log_list](#slow_log_list_struct) structure is documented below.

<a name="slow_log_list_struct"></a>
The `slow_log_list` block supports:

* `client_ip` - Indicates the IP address.

* `count` - Indicates the number of executions.

* `database` - Indicates the database that slow query logs belong to.

* `lock_time` - Indicates the average lock wait time.

* `node_id` - Indicates the ID of the instance node.

* `query_sample` - Indicates the execution syntax.

* `rows_examined` - Indicates the average number of scanned rows.

* `rows_sent` - Indicates the average number of result rows.

* `time` - Indicates the average execution time.

* `type` - Indicates the statement type.

* `users` - Indicates the name of the account.
