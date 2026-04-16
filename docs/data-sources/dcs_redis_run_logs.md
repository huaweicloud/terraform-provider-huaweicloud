---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_redis_run_logs"
description: |-
  Use this data source to query the list of Redis running logs within HuaweiCloud.
---

# huaweicloud_dcs_redis_run_logs

Use this data source to query the list of Redis running logs within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_redis_run_logs" "test" {
  instance_id = var.instance_id
  log_type    = "run"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the DCS instance is located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

* `log_type` - (Required, String) Specifies the type of log to query.
  Currently, only **run** (Redis running log) is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `file_list` - The list of running log files.  
  The [file_list](#dcs_redislogs_file) structure is documented below.

<a name="dcs_redislogs_file"></a>
The `file_list` block supports:

* `id` - The unique identifier of the log.

* `file_name` - The running log file name.

* `group_name` - The shard name.

* `replication_ip` - The IP address of the replica where the running log was collected.

* `status` - The status of getting the running log. Valid values are:
  + **succeed**: Success.
  + **failed**: Failed.

* `time` - The date when the running log was collected, in the format **yyyy-MM-dd**.

* `backup_id` - The ID of the log file.
