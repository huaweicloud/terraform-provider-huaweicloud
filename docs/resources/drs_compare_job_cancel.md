---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_compare_job_cancel"
description: |-
  Manages a resource to cancel a DRS compare job within HuaweiCloud.
---

# huaweicloud_drs_compare_job_cancel

Manages a resource to cancel a DRS compare job within HuaweiCloud.

-> This resource is a one-time action resource used to cancel a compare job. Deleting this resource will not
  undo the cancel operation, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "job_id" {} 
variable "compare_job_id" {}

resource "huaweicloud_drs_compare_job_cancel" "test" {
  job_id         = var.job_id 
  compare_job_id = var.compare_job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `job_id` - (Required, String, NonUpdatable) Specifies the DRS job ID.

* `compare_job_id` - (Required, String, NonUpdatable) Specifies the compare job ID to be canceled.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
