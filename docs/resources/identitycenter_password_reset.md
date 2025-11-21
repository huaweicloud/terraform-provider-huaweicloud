---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_password_reset"
description: |-
  Manages an Identity Center password reset resource within HuaweiCloud.
---

# huaweicloud_identitycenter_password_reset

Manages an Identity Center password reset resource within HuaweiCloud.

-> This resource is only a one-time action resource for operating the API.
  Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "user_id" {}
variable "identity_store_id" {}

resource "huaweicloud_identitycenter_password_reset" "test"{
  mode              = "EMAIL"
  user_id           = var.user_id
  identity_store_id = var.identity_store_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `identity_store_id` - (Required, String, NonUpdatable) Specifies the ID of the identity store.

* `user_id` - (Required, String, NonUpdatable) Specifies the ID of the user.

* `mode` - (Required, String, NonUpdatable) Specifies the reset mode. Value options: **EMAIL**, **OTP**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
