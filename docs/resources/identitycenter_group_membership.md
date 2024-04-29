---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_group_membership"
description: ""
---

# huaweicloud_identitycenter_group_membership

Manages an Identity Center group membership resource within HuaweiCloud.

## Example Usage

```hcl
variable "group_id" {}
variable "member_id" {}

data "huaweicloud_identitycenter_instance" "system" {}

resource "huaweicloud_identitycenter_group_membership" "test"{
  identity_store_id = data.huaweicloud_identitycenter_instance.system.identity_store_id
  group_id          = var.group_id
  member_id         = var.member_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `identity_store_id` - (Required, String, ForceNew) Specifies the ID of the identity store.

  Changing this parameter will create a new resource.

* `group_id` - (Required, String, ForceNew) Specifies the ID of the group.

  Changing this parameter will create a new resource.

* `member_id` - (Required, String, ForceNew) Specifies the ID of the user.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The Identity Center group membership can be imported using the `identity_store_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_identitycenter_group_membership.test <identity_store_id>/<id>
```
