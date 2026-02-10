---
subcategory: "Object Storage Migration Service (OMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_oms_sync_task_statistics"
description: |-
  Use this data source to query the synchronization task statistics.
---

# huaweicloud_oms_sync_task_statistics

Use this data source to query the synchronization task statistics.

## Example Usage

```hcl
variable "sync_task_id "{}
variable "data_type "{}
variable "start_time "{}
variable "end_time "{}

data "huaweicloud_oms_sync_task_statistics" "test" {
  sync_task_id = var.sync_task_id
  data_type    = var.data_type
  start_time   = var.start_time
  end_time     = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `sync_task_id` - (Required, String) Specifies the synchronization task ID.

* `data_type` - (Required, String) Specifies the statistical data type.  
  Use commas (,) to separate multiple data types.
  The valid values are as follows:
  + **REQUEST**: The number of objects requested for synchronization.
  + **SUCCESS**: The number of objects that are successfully synchronized.
  + **FAILURE**: The number of objects that fail to be synchronized.
  + **SKIP**: The number of objects that are skipped during synchronization.
  + **SIZE**: The size of successfully synchronized objects, in bytes.

* `start_time` - (Required, String) Specifies the start time for the query.
  The format is 13-digit timestamp in millisecond.

* `end_time` - (Required, String) Specifies the end time for the query.
  The format is 13-digit timestamp in millisecond.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `task_id` - The synchronization task ID.

* `statistic_time_type` - The interval for collecting statistics.
  + **FIVE_MINUTES**
  + **ONE_HOUR**

* `statistic_datas` - The statistics of the queried synchronization task.
  The [statistic_datas](#oms_sync_task_statistics_statistic_datas_struct) structure is documented below.

<a name="oms_sync_task_statistics_statistic_datas_struct"></a>
The `statistic_datas` block supports:

* `data_type` - The synchronization task ID.

* `data` - The source cloud service provider.
  The [data](#oms_sync_task_statistics_data_struct) structure is documented below.

<a name="oms_sync_task_statistics_data_struct"></a>
The `data` block supports:

* `time_stamp` - The statistics timestamp.

* `statistic_num` - The statistics number.
