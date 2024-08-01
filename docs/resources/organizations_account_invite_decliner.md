---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_account_invite_decliner"
description: |-
  Manages an Organizations account invite decliner resource within HuaweiCloud.
---

# huaweicloud_organizations_account_invite_decliner

Manages an Organizations account invite decliner resource within HuaweiCloud.

## Example Usage

```hcl
variable invitation_id {}

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
