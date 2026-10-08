---
subcategory: "Resource Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rc_resource_batch_add"
description: |-
  Use this resource to batch add resource to the resource group within HuaweiCloud.
---

# huaweicloud_rc_resource_batch_add

Use this resource to batch add resource to the resource group within HuaweiCloud.

-> This resource is a one-time action resource used to batch add resource. Deleting this resource will
not clear the corresponding request record, but will only remove the resource from the tfstate file.

## Example Usage

```hcl
variable "group_id" {}
variable "resource_ids" {
  type = list(string)
}

resource "huaweicloud_rc_resource_batch_add" "test" {
  group_id     = var.group_id
  resource_ids = var.resource_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the resources to be added are located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `group_id` - (Required, String, NonUpdatable) Specifies the ID of the resource group.

* `resource_ids` - (Required, List, NonUpdatable) Specifies the list of resource IDs to be added.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
