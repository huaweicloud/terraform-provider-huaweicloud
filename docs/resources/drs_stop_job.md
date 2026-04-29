---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_stop_job"
description: |-
  Manages a resource to stop a DRS job within HuaweiCloud.
---

# huaweicloud_drs_stop_job

Manages a resource to stop a DRS job within HuaweiCloud.

-> 1. This resource is a one-time action resource used to stop a DRS job. Deleting this resource will not restore the
  job or undo the stop action, but will only remove the resource information from the tf state file.
  <br/>2. You must specify an existing DRS job ID. If the job does not exist or has already ended,
  the operation may fail.<br/>3. The execution result of this operation is based on the value of the `status` field.

## Example Usage

```hcl
variable "job_id" {}

resource "huaweicloud_drs_stop_job" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `job_id` - (Required, String, NonUpdatable) Specifies the job ID.

* `is_force_stop` - (Optional, Bool, NonUpdatable) Specifies whether to force stop the job.
  The value can be **true** or **false**. Defaults to **false**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `job_id`.

* `name` - The name of the job.

* `status` - The status of the job.  
  The valid values are as follows:
  + **success**.
  + **failed**.
