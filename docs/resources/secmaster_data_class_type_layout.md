---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_data_class_type_layout"
description: |-
  Manages a resource to bind data class type with layout within HuaweiCloud.
---

# huaweicloud_secmaster_data_class_type_layout

Manages a resource to bind data class type with layout within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not unbind the layout from the data class
  type, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "workspace_id" {}
variable "dataclass_id" {}
variable "type_id" {}

resource "huaweicloud_secmaster_data_class_type_layout" "test" {
  workspace_id = var.workspace_id
  dataclass_id = var.dataclass_id
  type_id      = var.type_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID.

* `dataclass_id` - (Required, String, NonUpdatable) Specifies the unique ID of the data class.

* `type_id` - (Required, String, NonUpdatable) Specifies the type ID of the data class to which the layout will be bound.

* `layout_id` - (Optional, String, NonUpdatable) Specifies the ID of the layout to bind with the data class type.

* `layout_name` - (Optional, String, NonUpdatable) Specifies the name of the layout to bind with the data class type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
