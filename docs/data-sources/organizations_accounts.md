---
subcategory: "Organizations"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_organizations_accounts"
description: |-
  Use this data source to get the list of accounts in an organization.
---

# huaweicloud_organizations_accounts

Use this data source to get the list of accounts in an organization.

## Example Usage

```hcl
variable "parent_id" {}

data "huaweicloud_organizations_accounts" "test" {
  parent_id = var.parent_id
}
```

## Argument Reference

The following arguments are supported:

* `parent_id` - (Optional, String) Specifies the ID of root or organizational unit.

* `name` - (Optional, String) Specifies the name of the account.

* `with_register_contact_info` - (Optional, Bool) Whether to return email addresses and mobile
  numbers associated with the account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `accounts` - The list of accounts in an organization.
  The [accounts](#Organizations_Accounts) structure is documented below.

<a name="Organizations_Accounts"></a>
The `accounts` block supports:

* `id` - The ID of the account.

* `name` - The name of the account.

* `urn` - The uniform resource name of the account.

* `description` - The description of the account.

* `status` - The status of the account.

* `join_method` - How the account joined an organization.

* `joined_at` - The time when the account joined an organization.

* `mobile_phone` - The mobile phone number.

* `intl_number_prefix` - The prefix of a mobile phone number.

* `email` - The email address associated with the account.
