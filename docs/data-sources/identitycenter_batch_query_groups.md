---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_batch_query_groups"
description: |-
  Use this data source to get the Identity Center groups by group id.
---

# huaweicloud_identitycenter_batch_query_groups

Use this data source to get the Identity Center groups by group id.

## Example Usage

```hcl
variable "identity_store_id" {}
variable "group_ids" {}

data "huaweicloud_identitycenter_batch_query_groups" "test"{
  identity_store_id = var.identity_store_id
  group_ids         = var.group_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `identity_store_id` - (Required, String) Specifies the ID of the identity store that associated with IAM Identity
  Center.

* `group_ids` - (Required, List) Specifies the list of the group id.

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
