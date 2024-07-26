---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_studio_data_connections"
description: |-
  Use this data source to get the list of the Factory jobs within HuaweiCloud.
---
# huaweicloud_dataarts_factory_jobs

Use this data source to get the list of the Factory jobs within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_factory_jobs" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Optional, String) The ID of the workspace to which the jobs belong.
  If omitted, default workspace will be used.

* `name` - (Optional, String) Specified the job name to be queried. Fuzzy search is supported.

* `process_type` - (Optional, String) Specified the job type to be queried.
  If omitted, the default value is **BATCH**.  
  The valid values are as follows:
  + **REAL_TIME**: Real-time processing.
  + **BATCH**: Batch processing.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - All jobs that match the filter parameters.
  The [jobs](#factory_jobs) structure is documented below.

<a name="factory_jobs"></a>
The `jobs` block supports:

* `name` - The name of the job.

* `process_type` - The type of the job.

* `priority` - The priority of the job.
  + **0**: High priority.
  + **1**: Medium priority.
  + **2**: Low priority.

* `owner` - The owner of the job.

* `is_single_task_job` - Whether the job is single task.

* `directory` - The directory tree path of the job.

* `status` - The current status of the job.
  + **NORMAL**
  + **STOPPED**
  + **SCHEDULING**
  + **PAUSED**
  + **EXCEPTION**

* `start_time` - The start time of the job scheduling, in RFC3339 format.

* `end_time` - The end time of the job scheduling, in RFC3339 format.

* `created_by` - The creator of the job.

* `created_at` - The creation time of the job, in RFC3339 format.

* `updated_by` - The name of the user who last updated the job.

* `updated_at` - The latest update time of the job, in RFC3339 format.

* `last_instance_status` - The latest running status of the instance corresponding to the job.
  + **running**
  + **success**
  + **fail**
  + **running-exception**
  + **manual-stop**

* `last_instance_end_time` - The latest end time of the instance corresponding to the job, in RFC3339 format.
