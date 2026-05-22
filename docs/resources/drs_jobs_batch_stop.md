---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_jobs_batch_stop"
description: |-
  Manages a resource to batch stop DRS jobs within HuaweiCloud.
---

# huaweicloud_drs_jobs_batch_stop

Manages a resource to batch stop DRS jobs within HuaweiCloud.

-> 1. This resource is a one-time action resource used to batch stop DRS jobs. Deleting this resource will not restore
  the stopped jobs or undo the stop action, but will only remove the resource information from the tf state file.
  <br/>2. You must specify existing DRS job IDs. If a job does not exist or has already ended, the operation for that
  job may fail.<br/>3. The execution result of this operation is based on the `status` field in the `results` block.

## Example Usage

```hcl
variable "jobs" {
  type = list(object({
    job_id        = string
    is_force_stop = optional(bool, false)
  }))
}

resource "huaweicloud_drs_jobs_batch_stop" "test" {
  dynamic "jobs" {
    for_each = var.jobs

    content {
      job_id        = jobs.value.job_id
      is_force_stop = jobs.value.is_force_stop
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `jobs` - (Required, List, NonUpdatable) Specifies the request body for batch stopping tasks.
  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `job_id` - (Required, String, NonUpdatable) Specifies the ID of the DRS job.

* `is_force_stop` - (Optional, Bool, NonUpdatable) Specifies whether to force stop the job.
  The value can be **true** or **false**. Defaults to **false**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `results` - The response body for batch operation tasks.

  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `id` - The ID of the DRS job.

* `name` - The name of the DRS job.

* `status` - The operation result.
  The valid values are as follows:
  + **success**: The operation succeeded.
  + **failed**: The operation failed.

* `error_code` - The error code.

* `error_msg` - The error message.
