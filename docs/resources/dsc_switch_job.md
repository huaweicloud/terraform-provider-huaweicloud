---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_switch_job"
description: |-
  Manages a resource to switch the status of a DSC scan job within HuaweiCloud.
---

# huaweicloud_dsc_switch_job

Manages a resource to switch the status of a DSC scan job within HuaweiCloud.

-> This resource is a one-time action resource used to switch a DSC scan job status. Deleting this resource will not
  revert the job status change, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "job_id" {}

resource "huaweicloud_dsc_switch_job" "test" {
  job_id = var.job_id
  open   = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `job_id` - (Required, String, NonUpdatable) Specifies the scan job ID.

* `open` - (Optional, Bool) Specifies the job switch status. Defaults to **false**.
  The valid values are as follows:
  + **true**: Enable the scan job.
  + **false**: Disable the scan job.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `msg` - The returned message, which describes the operation result or status information.

* `status` - The returned status, indicating whether the operation is successful.
