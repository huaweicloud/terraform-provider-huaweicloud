---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_resource_share_accepter"
description: |-
  Manages a RAM resource share accepter resource within HuaweiCloud.
---

# huaweicloud_ram_resource_share_accepter

Manages a RAM resource share accepter resource within HuaweiCloud.

## Example Usage

```hcl
variable "invitation_id" {}
variable "action" {}

resource "huaweicloud_ram_resource_share_accepter" "test" {
  resource_share_invitation_id = var.invitation_id
  action                       = var.action
}
```

## Argument Reference

The following arguments are supported:

* `resource_share_invitation_id` - (Required, String, ForceNew) Specifies the ID of the resource share invitation.

* `action` - (Required, String, ForceNew) Specifies the action to be taken for resource share invitation.
  The value can be **accept** or **reject**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `resource_share_invitation_id`.
