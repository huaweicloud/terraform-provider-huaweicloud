---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_slow_logs"
description: |-
  Use this data source to get the list of RDS slow logs.
---

# huaweicloud_rds_slow_logs

Use this data source to get the list of RDS slow logs.

## Example Usage

```hcl
variable "instance_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_rds_slow_logs" "test" {
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

* `type` - (Optional, String) Specifies the statement type. Value options: **INSERT**, **UPDATE**, **SELECT**,
  **DELETE**, **CREATE**.

* `database` - (Optional, String) Specifies the name of the database.

* `users` - (Optional, String) Specifies the name of the account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `slow_logs` - Indicates the list of the slow logs.

  The [slow_logs](#slow_logs_struct) structure is documented below.

<a name="slow_logs_struct"></a>
The `slow_logs` block supports:

* `count` - Indicates the number of execution times.

* `time` - Indicates the execution time.

* `lock_time` - Indicates the wait lock time.

* `rows_sent` - Indicates the number of result lines.

* `rows_examined` - Indicates the number of rows scanned.

* `database` - Indicates the name of the database.

* `users` - Indicates the name of the account.

* `query_sample` - Indicates the execution syntax.

* `type` - Indicates the statement type.

* `start_time` - Indicates the start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `client_ip` - Indicates the IP address of the client.
