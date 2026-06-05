---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_alarm_statistics"
description: |-
  Use this data source to query GaussDB instance alarm statistics within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_alarm_statistics

Use this data source to query GaussDB instance alarm statistics within HuaweiCloud.

## Example Usage

```hcl
variable "start_time" {}

variable "top_num" {}

data "huaweicloud_gaussdb_instance_alarm_statistics" "test" {
  start_time = var.start_time
  top_num    = var.top_num
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `start_time` - (Required, String) Specifies the start time for querying alarm statistics.
  The value is in the format of **yyyy-mm-ddThh:mm:ss+0000**.

* `top_num` - (Required, Int) Specifies the number of instances with the most alarms to return.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ring_percentage` - The ring percentage of alarms.

* `instance_alarm_level_statistics` - The list of alarm level statistics by instance.
  The [instance_alarm_level_statistics](#instance_alarm_level_statistics_struct) structure is documented below.

* `total_alarm_level_statistics` - The list of total alarm level statistics.
  The [total_alarm_level_statistics](#total_alarm_level_statistics_struct) structure is documented below.

<a name="instance_alarm_level_statistics_struct"></a>
The `instance_alarm_level_statistics` block supports:

* `instance_id` - The ID of the GaussDB instance.

* `instance_name` - The name of the GaussDB instance.

* `total_count` - The total number of alarms for the instance.

* `alarm_level_statistics` - The list of alarm level statistics for the instance.
  The [alarm_level_statistics](#alarm_level_statistics_struct) structure is documented below.

<a name="total_alarm_level_statistics_struct"></a>
The `total_alarm_level_statistics` block supports:

* `count` - The count of alarms at this level.

* `level_name` - The name of the alarm level.

<a name="alarm_level_statistics_struct"></a>
The `alarm_level_statistics` block supports:

* `count` - The count of alarms at this level.

* `level_name` - The name of the alarm level.
