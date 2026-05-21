---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_batch_retry_task"
description: |-
  Manages a resource to batch retry DRS tasks within HuaweiCloud.
---

# huaweicloud_drs_batch_retry_task

Manages a resource to batch retry DRS tasks within HuaweiCloud.

-> 1. This resource is a one-time action resource used to retry DRS tasks. Deleting this resource will not
  undo the retry operation, but will only remove the resource information from the tf state file.
  <br/>2. Tasks in failed status can be retried.
  <br/>3. A successful API call does not guarantee the operation success; please check the task status.

## Example Usage

```hcl
variable "jobs" {
  type = list(object({
    job_id          = string
    is_sync_re_edit = string
  }))
}

resource "huaweicloud_drs_batch_retry_task" "test" {
  dynamic "jobs" {
    for_each = var.jobs

    content {
      job_id          = jobs.value.job_id
      is_sync_re_edit = jobs.value.is_sync_re_edit
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `jobs` - (Required, List, NonUpdatable) Specifies the list of jobs to retry.
  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `job_id` - (Required, String, NonUpdatable) Specifies the job ID.

* `is_sync_re_edit` - (Optional, String, NonUpdatable) Specifies whether the task is started after re-editing.
  The valid values are as follows:
  + **true**: The task is started after re-editing.
  + **false**: The task is not started after re-editing.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `results` - The results of retrying tasks.
  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `id` - The job ID.

* `status` - The retry operation result. Valid values are **success** and **failed**.

* `error_code` - The error code.

* `error_msg` - The error message.
