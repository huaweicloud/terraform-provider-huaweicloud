---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_maintenance_batch_manage"
description: |-
  Manages a Workspace desktop maintenance mode resource within HuaweiCloud.
---

# huaweicloud_workspace_desktop_maintenance_batch_manage

Manages a Workspace desktop maintenance mode resource within HuaweiCloud.

-> This resource is a one-time action resource for batch managing desktops' maintenance mode. Deleting this resource
   will not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "desktop_ids" {
  type = list(string)
}

resource "huaweicloud_workspace_desktop_maintenance_batch_manage" "test" {
  desktop_ids         = var.desktop_ids
  in_maintenance_mode = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the desktops are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `desktop_ids` - (Required, List, NonUpdatable) Specifies the list of desktop IDs to set maintenance mode.
  
* `in_maintenance_mode` - (Required, Bool, NonUpdatable) Specifies whether to enter or exit maintenance mode.
  + **true**: Enter maintenance mode
  + **false**: Exit maintenance mode

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
