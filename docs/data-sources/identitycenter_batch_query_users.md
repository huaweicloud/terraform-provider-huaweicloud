---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_batch_query_users"
description: |-
  Use this data source to get the Identity Center users by user id.
---

# huaweicloud_identitycenter_batch_query_users

Use this data source to get the Identity Center users by user id.

## Example Usage

```hcl
variable "identity_store_id" {}
variable "user_ids" {}
 
data "huaweicloud_identitycenter_batch_query_users" "test"{
  identity_store_id = var.identity_store_id
  user_ids          = var.user_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `identity_store_id` - (Required, String) Specifies the ID of the identity store that associated with IAM Identity
  Center.

* `user_ids` - (Required, List) Specifies the list of the user id.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `users` - Indicates the list of IdentityCenter user.
  The [users](#IdentityCenterUsers_User) structure is documented below.

<a name="IdentityCenterUsers_User"></a>
The `users` block supports:

* `id` - Indicates the ID of the user.

* `user_name` - Indicates the name of the user.

* `family_name` - Indicates the family name of the user.

* `given_name` - Indicates the given name of the user.

* `display_name` - Indicates the display name of the user.

* `email` - Indicates the email of the user.

* `phone_number` - The phone number of the user.

* `user_type` - The type of the user.

* `title` - The title of the user.

* `addresses` - The addresses information of the user.
  The [addresses](#addresses_struct) structure is documented below.

* `enterprise` - The enterprise information of the user.
  The [enterprise](#enterprise_struct) structure is documented below.

* `created_at` - The creation time of the user.

* `created_by` - The creator of the user.

* `updated_at` - The update time of the user.

* `updated_by` - The updater of the user.

* `email_verified` - Whether the email is verified.

* `enabled` - Whether the user is enabled.

<a name="addresses_struct"></a>
The `addresses` block supports:

* `country` - The country of the user.

* `formatted` - A string containing a formatted version of the address to be displayed.

* `locality` - The locality of the user.

* `postal_code` - The postal code of the user.

* `region` - The region of the user.

* `street_address` - The street address of the user.

<a name="enterprise_struct"></a>
The `enterprise` block supports:

* `cost_center` - The cost center of the enterprise.

* `department` - The department of the enterprise.

* `division` - The division of the enterprise.

* `employee_number` - The employee number of the enterprise.

* `organization` - The organization of the enterprise.

* `manager` - The manager of the enterprise.
