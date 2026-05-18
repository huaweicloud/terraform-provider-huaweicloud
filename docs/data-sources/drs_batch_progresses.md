---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_batch_progresses"
description: |-
  Use this data source to get the progresses list for specified DRS jobs within HuaweiCloud.
---

# huaweicloud_drs_batch_progresses

Use this data source to get the progresses list for specified DRS jobs within HuaweiCloud.

## Example Usage

```hcl
variable "job_ids" { 
  type = list(string) 
}

data "huaweicloud_drs_batch_progresses" "test" { 
  job_ids = var.job_ids 
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_ids` - (Required, List) Specifies the list of DRS job IDs to query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `results` - The progresses list.

  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `job_id` - The ID of the DRS job.

* `progress` - The migration percentage.

* `incre_trans_delay` - The incremental migration delay in seconds.

* `incre_trans_delay_millis` - The incremental migration delay in milliseconds.

* `task_mode` - The task mode.
  The valid values are as follows:
  + **FULL_TRANS**: Full migration.
  + **INCR_TRANS**: Incremental migration.
  + **FULL_INCR_TRANS**: Full + incremental migration.

* `transfer_status` - The task status.
  The valid values are as follows:
  + **CREATING**: Creating.
  + **CREATE_FAILED**: Creation failed.
  + **CONFIGURATION**: Configuring.
  + **WAITING_FOR_START**: Waiting to start.
  + **RELEASE_RESOURCE_COMPLETE**: Completed.
  + **DELETED**: Deleted.
  + **INCRE_TRANSFER_STARTED**: Incremental migration in progress.
  + **INCRE_TRANSFER_FAILED**: Incremental migration failed.
  + **FULL_TRANSFER_STARTED**: Full migration in progress.
  + **FULL_TRANSFER_COMPLETE**: Full migration completed.
  + **PAUSING**: Pausing.
  + **FULL_TRANSFER_FAILED**: Full migration failed.

* `process_time` - The migration time, timestamp.

* `remaining_time` - The estimated remaining time.

* `progress_map` - The data, structure, and index migration progress information.

  The [progress_map](#progress_map_struct) structure is documented below.

* `error_code` - The error code.

* `error_msg` - The error message.

<a name="progress_map_struct"></a>
The `progress_map` block supports:

* `key` - The key of the progress map.

* `completed` - The completion progress.

* `remaining_time` - The estimated remaining time.
