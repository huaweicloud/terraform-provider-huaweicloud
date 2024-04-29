---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_groups"
description: ""
---

# huaweicloud_identitycenter_groups

Use this data source to get the Identity Center groups.

## Example Usage

```hcl
data "huaweicloud_identitycenter_instance" "system" {}

data "huaweicloud_identitycenter_groups" "test"{
  identity_store_id = data.huaweicloud_identitycenter_instance.system.identity_store_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `identity_store_id` - (Required, String) Specifies the ID of the identity store that associated with IAM Identity
  Center.

* `name` - (Optional, String) Specifies the name of the group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `groups` - Indicates the list of IdentityCenter group.
  The [groups](#IdentityCenterGroups_Group) structure is documented below.

<a name="IdentityCenterGroups_Group"></a>
The `groups` block supports:

* `id` - Indicates the ID of the group.

* `name` - Indicates the name of the group.

* `description` - Indicates the description of the group.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the last update time.
