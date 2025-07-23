---
subcategory: "Storage Disaster Recovery Service (SDRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sdrs_delete_all_group_failure_jobs"
description: |-
  Using this resource to delete all failure jobs of all protection groups in SDRS within HuaweiCloud.
---

# huaweicloud_sdrs_delete_all_group_failure_jobs

Using this resource to delete all failure jobs of all protection groups in SDRS within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not recover the deleted failure jobs,
but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
resource "huaweicloud_sdrs_delete_all_group_failure_jobs" "test" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
