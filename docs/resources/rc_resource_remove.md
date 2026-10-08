---
subcategory: "Resource Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rc_resource_remove"
description: |-
  Use this resource to remove resource from the resource group within HuaweiCloud.
---

# huaweicloud_rc_resource_remove

Use this resource to remove resource from the resource group within HuaweiCloud.

-> This resource is a one-time action resource used to remove resource. Deleting this resource will
not clear the corresponding request record, but will only remove the resource from the tfstate file.

## Example Usage

```hcl
variable "group_id" {}
variable "resource_id" {}

resource "huaweicloud_rc_resource_remove" "test" {
  group_id    = var.group_id
  resource_id = var.resource_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the resource to be removed is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `group_id` - (Required, String, NonUpdatable) Specifies the ID of the resource group.

* `resource_id` - (Required, String, NonUpdatable) Specifies the ID of the resource to be removed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
