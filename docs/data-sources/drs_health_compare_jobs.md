---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_health_compare_jobs"
description: |-
  Use this data source to get the list of health compare jobs for a specified DRS task within a region.
---

# huaweicloud_drs_health_compare_jobs

Use this data source to get the list of health compare jobs for a specified DRS task within a region.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_health_compare_jobs" "all_jobs" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS task.

* `status` - (Optional, String) Specifies the status of the compare jobs. If omitted, all jobs will be queried.
  The valid values are as follows:
  + **WAITING_FOR_RUNNING**: Waiting to start.
  + **RUNNING**: Running.
  + **SUCCESSFUL**: Succeeded.
  + **FAILED**: Failed.
  + **CANCELLED**: Cancelled.
  + **TIMEOUT_INTERRUPT**: Timeout interrupted.
  + **FULL_DOING**: Full comparison in progress.
  + **INCRE_DOING**: Incremental comparison in progress.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `compare_jobs` - The list of health compare jobs.

  The [compare_jobs](#compare_jobs_struct) structure is documented below.

<a name="compare_jobs_struct"></a>
The `compare_jobs` block supports:

* `id` - The compare job ID.

* `type` - The comparison type.
  The valid values are as follows:
  + **object_comparison**: Object comparison.
  + **lines**: Row comparison.
  + **account**: User comparison.

* `start_time` - The start time of the comparison.

* `end_time` - The end time of the comparison.

* `status` - The status of the comparison.
  The valid values are as follows:
  + **WAITING_FOR_RUNNING**: Waiting to start.
  + **RUNNING**: Running.
  + **SUCCESSFUL**: Succeeded.
  + **FAILED**: Failed.
  + **CANCELLED**: Cancelled.
  + **TIMEOUT_INTERRUPT**: Timeout interrupted.
  + **FULL_DOING**: Full comparison in progress.
  + **INCRE_DOING**: Incremental comparison in progress.

* `compute_type` - The compute resource type used for the comparison.

* `database_info` - The database information involved in the comparison.

  The [database_info](#database_info_struct) structure is documented below.

<a name="database_info_struct"></a>
The `database_info` block supports:

* `service_database` - The service database information.

* `disaster_recovery_database` - The disaster recovery database information.
