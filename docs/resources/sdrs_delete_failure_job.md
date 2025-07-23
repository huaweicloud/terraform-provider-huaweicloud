---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_delete_failure_job"
description: |-
  Using this resource to delete a failure job in SDRS within HuaweiCloud.
---

# huaweicloud_sdrs_delete_failure_job

Using this resource to delete a failure job in SDRS within HuaweiCloud.

-> This is a one-time action resource to delete a failure job. Deleting this resource will
not change the current configuration, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "failure_job_id" {}

resource "huaweicloud_sdrs_delete_failure_job" "test" {
  failure_job_id = var.failure_job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `failure_job_id` - (Required, String, NonUpdatable) Specifies the ID of the failure job to delete.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
