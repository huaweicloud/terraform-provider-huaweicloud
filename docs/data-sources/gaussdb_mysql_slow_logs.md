---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_slow_logs"
description: |-
  Use this data source to get the list of GaussDB MySQL slow logs.
---

# huaweicloud_gaussdb_mysql_slow_logs

Use this data source to get the list of GaussDB MySQL slow logs.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_gaussdb_mysql_slow_logs" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  start_time  = var.start_time
  end_time    = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB MySQL instance.

* `node_id` - (Required, String) Specifies the ID of the instance node.

* `start_time` - (Required, String) Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - (Required, String) Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `operate_type` - (Optional, String) Specifies the SQL statement type. Value options: **INSERT**, **UPDATE**, **SELECT**,
  **DELETE**, **CREATE**, **ALTER**, **DROP**.

* `database` - (Optional, String) Specifies the name of the database.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `slow_log_list` - Indicates the list of the slow logs.

  The [slow_log_list](#slow_log_list_struct) structure is documented below.

<a name="slow_log_list_struct"></a>
The `slow_log_list` block supports:

* `count` - Indicates the ID of the instance node.

* `count` - Indicates the number of executions.

* `time` - Indicates the execution time.

* `lock_time` - Indicates the lock wait time.

* `rows_sent` - Indicates the number of sent rows.

* `rows_examined` - Indicates the number of scanned rows.

* `database` - Indicates the database that slow query logs belong to.

* `users` - Indicates the name of the account.

* `query_sample` - Indicates the execution syntax.

* `type` - Indicates the statement type.

* `start_time` - Indicates the start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `client_ip` - Indicates the IP address of the client.
