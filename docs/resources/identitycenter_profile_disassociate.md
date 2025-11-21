---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_profile_disassociate"
description: |-
  Manages an Identity Center profile disassociate resource within HuaweiCloud.
---

# huaweicloud_identitycenter_profile_disassociate

Manages an Identity Center profile disassociate resource within HuaweiCloud.

-> This resource is only a one-time action resource for operating the API.
  Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "identity_store_id" {}
variable "accessor_id" {}

resource "huaweicloud_identitycenter_profile_disassociate" "test"{
  instance_id       = var.instance_id
  identity_store_id = var.identity_store_id
  accessor_type     = "USER"
  accessor_id       = var.accessor_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Identity Center instance.

* `identity_store_id` - (Required, String, NonUpdatable) Specifies the ID of the identity store.

* `accessor_id` - (Required, String, NonUpdatable) Specifies the ID of the accessor.

* `accessor_type` - (Required, String, NonUpdatable) Specifies the type of the accessor. Value options: **USER**, **GROUP**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
