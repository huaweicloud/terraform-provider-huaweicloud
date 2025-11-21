---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_groups"
description: |-
  Use this data source to get the list of group in the Identity and Access Management V5 service.
---

# huaweicloud_identityv5_groups

Use this data source to get the list of group in the Identity and Access Management V5 service.

## Example Usage

```hcl
data "huaweicloud_identityv5_groups" "groups" {}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Optional, String) Specifies the ID of the IAM user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `groups` - The list of groups. Each element contains the following attributes:
  The [groups](#Identity_Group) structure is documented below.

<a name="Identity_Group"></a>
The `groups` block supports:

* `group_name` - Indicates the name of the group.

* `urn` - Indicates the Uniform Resource Name (URN) of the group.

* `created_at` - Indicates the time when the group was created.

* `description` - Indicates the description of the group.

* `group_id` - Indicates the ID of the group.
