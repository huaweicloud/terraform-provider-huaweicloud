---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_htap_starrocks_error_logs"
description: |-
  Use this data source to query the list of error logs of TaurusDB HTAP StarRocks instance within HuaweiCloud.
---

# huaweicloud_taurusdb_htap_starrocks_error_logs

Use this data source to query the list of error logs of TaurusDB HTAP StarRocks instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_taurusdb_htap_starrocks_error_logs" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  start_time  = var.start_time
  end_time    = var.end_time
  level       = "ALL"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the StarRocks HTAP instance.

* `node_id` - (Required, String) Specifies the ID of the instance node.

* `start_time` - (Required, String) Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - (Required, String) Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `level` - (Required, String) Specifies the log level.
  Valid values are as follows:
  + **ALL**
  + **INFO**
  + **WARNING**
  + **ERROR**
  + **FATAL**
  + **PANIC**
  + **NOTE**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `error_log_list` - Indicates the list of the error logs.

  The [error_log_list](#error_log_list_struct) structure is documented below.

<a name="error_log_list_struct"></a>
The `error_log_list` block supports:

* `node_id` - Indicates the ID of the instance node.

* `time` - Indicates the log time.

* `level` - Indicates the log level.

* `content` - Indicates the log content.
