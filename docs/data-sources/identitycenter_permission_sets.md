---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_permission_sets"
description: |-
  Use this data source to get the Identity Center permission sets.
---

# huaweicloud_identitycenter_permission_sets

Use this data source to get the Identity Center permission sets.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_identitycenter_permission_sets" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of an IAM Identity Center instance.

* `permission_set_id` - (Optional, String) Specifies the ID of a permission set.

* `name` - (Optional, String) Specifies the name of a permission set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `permission_sets` - The permission set list.

  The [permission_sets](#permission_sets_struct) structure is documented below.

<a name="permission_sets_struct"></a>
The `permission_sets` block supports:

* `description` - The description of a permission set.

* `name` - The name of a permission set.

* `permission_set_id` - The ID of a permission set.

* `relay_state` - The redirection of users within an application during the federated authentication.

* `session_duration` - The length of time that the application user sessions are valid.

* `permission_urn` - The URN of a permission set.

* `created_at` - The time when a permission set is created.
