---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_batch_delete_obs"
description: |-
  Manages a resource to batch delete OBS assets within HuaweiCloud.
---

# huaweicloud_dsc_batch_delete_obs

Manages a resource to batch delete OBS assets within HuaweiCloud.

-> This resource is a one-time action resource used to batch delete DSC OBS assets. Deleting this resource will not
  restore the deleted OBS assets or undo the delete action, but will only remove the resource information from
  the tf state file.

## Example Usage

```hcl
variable "obs_ids" {
  type = list(string)
}

resource "huaweicloud_dsc_batch_delete_obs" "test" {
  obs_ids = var.obs_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `obs_ids` - (Required, List, NonUpdatable) Specifies the list of OBS bucket IDs to delete.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `msg` - The returned message describing the operation result or error information.

* `status` - The returned status, such as '200' or '400'.
