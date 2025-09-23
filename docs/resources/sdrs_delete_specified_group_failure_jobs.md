---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_delete_specified_group_failure_jobs"
description: |-
  Using this resource to delete all failure jobs of a specified protection group in SDRS within HuaweiCloud.
---

# huaweicloud_sdrs_delete_specified_group_failure_jobs

Using this resource to delete all failure jobs of a specified protection group in SDRS within HuaweiCloud.

-> This is a one-time action resource to delete all failure jobs from a protected group. Deleting this
resource will not change the current configuration, but will only remove the resource information from the
tfstate file.

## Example Usage

```hcl
variable "server_group_id" {}

resource "huaweicloud_sdrs_delete_specified_group_failure_jobs" "test" {
  server_group_id = var.server_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `server_group_id` - (Required, String, NonUpdatable) Specifies the ID of the protected group to delete all failure
  jobs from.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource.
