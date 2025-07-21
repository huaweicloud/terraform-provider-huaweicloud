---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_sysprep"
description: |-
  Use this data source to get the sysprep information of a Workspace desktop within HuaweiCloud.
---

# huaweicloud_workspace_desktop_sysprep

Use this data source to get the sysprep information of a Workspace desktop within HuaweiCloud.

## Example Usage

```hcl
variable "desktop_id" {}

data "huaweicloud_workspace_desktop_sysprep" "test" {
  desktop_id = var.desktop_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the desktop sysprep is located.  
  If omitted, the provider-level region will be used.

* `desktop_id` - (Required, String) Specifies the ID of the desktop to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sysprep_info` - The sysprep information of the desktop.  
  The [sysprep_info](#workspace_desktop_sysprep_info) structure is documented below.

<a name="workspace_desktop_sysprep_info"></a>
The `sysprep_info` block supports:

* `sysprep_version` - The sysprep version of the desktop.

* `support_create_image` - Whether the desktop supports creating image.
