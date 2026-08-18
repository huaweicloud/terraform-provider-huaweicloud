---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_groups"
description: |-
  Use this data source to get the list of the IAM user groups.
---

# huaweicloud_identity_groups

Use this data source to get the list of the IAM user groups.

## Example Usage

### Query all user groups

```hcl
data "huaweicloud_identity_groups" "test" {}
```

### Query user groups by name

```hcl
variable "group_name" {}

data "huaweicloud_identity_groups" "test" {
  name = var.group_name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the name of the user group.

* `domain_id` - (Optional, String) Specifies the ID of the account to which the user groups belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - The list of user groups that match the filter parameters.  
  The [groups](#identity_user_groups) structure is documented below.

<a name="identity_user_groups"></a>
The `groups` block supports:

* `id` - The ID of the user group.

* `name` - The name of the user group.

* `description` - The description of the user group.

* `domain_id` - The ID of the account to which the user group belongs.

* `created_at` - The creation time of the user group, in RFC3339 format.
