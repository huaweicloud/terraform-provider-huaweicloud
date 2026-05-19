---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_slow_logs"
description: |-
  Use this data source to query slow logs of a DCS instance within HuaweiCloud.
---

# huaweicloud_dcs_slow_logs

Use this data source to query slow logs of a DCS instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dcs_slow_logs" "test" {
  instance_id = var.instance_id
  start_time  = "1598803200000"
  end_time    = "1599494399000"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the slow logs. If omitted, the provider-level
  region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DCS instance.

* `start_time` - (Required, String) Specifies the query start time, which is the Unix timestamp of UTC time. For
  example: **1598803200000**.

* `end_time` - (Required, String) Specifies the query end time, which is the Unix timestamp of UTC time. For example: *
  *1599494399000**.

* `sort_key` - (Optional, String) Specifies the sorting keyword.
  Valid values are **start_time** and **duration**. The default value is **start_time**.

* `sort_dir` - (Optional, String) Specifies the sorting direction.
  Valid values are **desc** (descending) and **asc** (ascending). The default value is **desc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `slowlogs` - The list of slow logs.
  The [slowlogs](#slowlogs_struct) structure is documented below.

<a name="slowlogs_struct"></a>
The `slowlogs` block supports:

* `id` - The unique identifier of the slow log.

* `command` - The slow command.

* `start_time` - The execution start time, in the format of "2020-06-19T07:06:07Z".

* `duration` - The execution duration in μs.

* `shard_name` - The shard name where the slow command is located.

* `database_id` - The database ID.

* `username` - The account name that operated the slow log.

* `node_role` - The node type.

* `client_ip` - The client IP address.
