---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_account_invite_decliner"
description: |-
  Use this resource to decline the received Organizations account invitation within HuaweiCloud.
---

# huaweicloud_organizations_account_invite_decliner

Use this resource to decline the received Organizations account invitation within HuaweiCloud.

-> This resource is only a one-time action resource for declining the received account invitation. Deleting this
   resource will not decline the received account invitation, but will only remove the resource information from the
   tfstate file.

## Example Usage

```hcl
variable "invitation_id" {}

resource "huaweicloud_organizations_account_invite_decliner" "test"{
  invitation_id = var.invitation_id
}
```

## Argument Reference

The following arguments are supported:

* `invitation_id` - (Required, String, ForceNew) Specifies the unique ID of an invitation (handshake).

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
