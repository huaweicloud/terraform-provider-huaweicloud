---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_volume_batch_delete"
description: |-
  Use this resource to batch delete desktop data volumes within HuaweiCloud.
---

# huaweicloud_workspace_desktop_volume_batch_delete

Use this resource to batch delete desktop data volumes within HuaweiCloud.

-> This resource is only a one-time action resource for batch deleting desktop data volumes. Deleting this
  resource will not clear the corresponding request record, but will only remove the resource information from the
  tfstate file.

## Example Usage

```hcl
variable "desktop_id" {}
variable "volume_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_desktop_volume_batch_delete" "test" {
  desktop_id = var.desktop_id
  volume_ids = var.volume_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the desktop is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `desktop_id` - (Required, String, NonUpdatable) Specifies the ID of the desktop to which the data volumes to
  be deleted belongs.

* `volume_ids` - (Required, List, NonUpdatable) Specifies the list of desktop data volume IDs to be deleted.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is `10` minutes.
