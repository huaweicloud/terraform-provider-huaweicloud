---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_groups"
description: |-
  Use this data source to get the list of the user groups within HuaweiCloud.
---

# huaweicloud_identityv5_groups

Use this data source to get the list of the user groups within HuaweiCloud.

## Example Usage

### Query all user groups

```hcl
data "huaweicloud_identityv5_groups" "test" {}
```

### Query user group by user ID

```hcl
variable "user_id" {}

data "huaweicloud_identityv5_groups" "test" {
  user_id = var.user_id
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Optional, String) Specifies the ID of the user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - The list of user groups that match the filter parameters.  
  The [groups](#v5_user_groups) structure is documented below.

<a name="v5_user_groups"></a>
The `groups` block supports:

* `group_id` - The ID of the user group.

* `group_name` - The name of the user group.

* `urn` - The uniform resource name of the user group.

* `description` - The description of the user group.

* `created_at` - The creation time of the user group.
