---
subcategory: "Resource Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rc_resource_group"
description: |-
  Manages a resource group within HuaweiCloud Resource Center.
---

# huaweicloud_rc_resource_group

Manages a resource group within HuaweiCloud Resource Center.

## Example Usage

```hcl
variable "group_name" {}
variable "group_description" {}

resource "huaweicloud_rc_resource_group" "test" {
  group_name  = var.group_name
  description = var.group_description
}
```

## Argument Reference

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `group_name` - (Required, String, ForceNew) Specifies the resource group name.

* `description` - (Optional, String) Specifies the resource group description.

## Attribute Reference

* `id` - The resource ID.

* `domain_id` - The domain ID of the resource group.

* `create_time` - The creation time of the resource group.

* `update_time` - The update time of the resource group.

## Import

This resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_rc_resource_group.test <id>
```
