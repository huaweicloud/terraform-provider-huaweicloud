---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_organization_accounts"
description: |-
  Use this data source to get the list of AOM organization accounts.
---

# huaweicloud_aom_organization_accounts

Use this data source to get the list of AOM organization accounts.

## Example Usage

```hcl
data "huaweicloud_aom_organization_accounts" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `accounts` - Indicates the accounts list.
  The [accounts](#accounts_struct) structure is documented below.

<a name="accounts_struct"></a>
The `accounts` block supports:

* `id` - Indicates the account ID.

* `name` - Indicates the account name.

* `urn` - Indicates the uniform resource name of the account.

* `join_method` - Indicates the method how the account joined in the organization.

* `joined_at` - Indicates the time when the account joined in the organization.
