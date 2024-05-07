---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_error_logs"
description: |-
  Use this data source to get the list of RDS error logs.
---

# huaweicloud_rds_error_logs

Use this data source to get the list of RDS error logs.

## Example Usage

```hcl
variable "instance_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_rds_error_logs" "test" {
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

* `level` - (Optional, String) Specifies the log level. Value options: **ALL**, **INFO**, **LOG**, **WARNING**,
  **ERROR**, **FATAL**, **PANIC**, **NOTE**. Defaults to **ALL**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `error_logs` - Indicates the list of the error logs.

  The [error_logs](#error_logs_struct) structure is documented below.

<a name="error_logs_struct"></a>
The `error_logs` block supports:

* `time` - Indicates the date and time of the error log in the **yyyy-mm-ddThh:mm:ssZ** format.

* `level` - Indicates the error log level.

* `content` - Indicates the error log content.
