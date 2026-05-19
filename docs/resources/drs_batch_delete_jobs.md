---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_batch_delete_jobs"
description: |-
  Manages a resource to batch delete DRS jobs within HuaweiCloud.
---

# huaweicloud_drs_batch_delete_jobs

Manages a resource to batch delete DRS jobs within HuaweiCloud.

-> 1. This resource is a one-time action resource used to batch delete DRS jobs. Deleting this resource
  will not restore the deleted jobs or undo the delete action, but will only remove the resource information from
  the tf state file.<br/>2. Before deleting jobs, please ensure that the jobs are in the finished status.

## Example Usage

```hcl
variable "job_ids" { 
  type = list(string)
}

resource "huaweicloud_drs_batch_delete_jobs" "test" { 
  jobs = var.job_ids 
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `jobs` - (Required, List, NonUpdatable) Specifies the list of job IDs to delete.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `results` - The results of batch deleting jobs.
  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `id` - The job ID.

* `name` - The job name.

* `status` - The delete operation result.
