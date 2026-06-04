---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_batch_set_definer"
description: |-
  Manages a resource to batch set DRS job definer within HuaweiCloud.
---

# huaweicloud_drs_batch_set_definer

Manages a resource to batch set DRS job definer within HuaweiCloud.

-> 1. This resource is a one-time action resource used to batch set DRS job definer. Deleting this resource will not
  restore the definer setting or undo the set action, but will only remove the resource information from the
  tf state file.
  <br/>2. You must specify existing DRS job IDs. The job status must be **CONFIGURATION** before calling this API.
  If a job does not exist or is not in the correct status, the operation may fail.<br/>3. The execution result of this
  operation is based on the `status` field in the `results` block.

## Example Usage

```hcl
variable "jobs" {
  type = list(object({
    job_id          = string
    replace_definer = bool
  }))
}

resource "huaweicloud_drs_batch_set_definer" "test" {
  dynamic "jobs" {
    for_each = var.jobs

    content {
      job_id          = jobs.value.job_id
      replace_definer = jobs.value.replace_definer
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `jobs` - (Required, List, NonUpdatable) Specifies the batch set replace definer request list.

  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `job_id` - (Required, String, NonUpdatable) Specifies the job ID.

* `replace_definer` - (Required, Bool, NonUpdatable) Specifies whether to use the target database user to replace the
  definer.  
  The valid values are as follows:
  + **true**: After migration, the definer of all source database objects will be migrated to the target user. Other
    users need to be authorized to have database object permissions.
  + **false**: After migration, the definer definition of source database objects will remain unchanged. This option
    requires migrating all source database users together with the user permission migration function to keep the
    source database permission system unchanged.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `results` - The batch modify task return list.

  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `id` - The job ID.

* `status` - The status of the operation.  
  The valid values are as follows:
  + **success**: The operation succeeded.
  + **failed**: The operation failed.

* `error_code` - The error code.

* `error_msg` - The error message.
