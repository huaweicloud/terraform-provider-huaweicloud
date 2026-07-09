---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_stop_scan_job"
description: |-
  Manages a resource to stop a DSC scan job within HuaweiCloud.
---

# huaweicloud_dsc_stop_scan_job

Manages a resource to stop a DSC scan job within HuaweiCloud.

-> This resource is a one-time action resource used to stop a DSC scan job. Deleting this resource will not
  restart the scan job, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "job_id" {}

resource "huaweicloud_dsc_stop_scan_job" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `job_id` - (Required, String, NonUpdatable) Specifies the scan job ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `msg` - The returned message, which describes the operation result or status information.

* `status` - The returned status, indicating whether the operation is successful.
