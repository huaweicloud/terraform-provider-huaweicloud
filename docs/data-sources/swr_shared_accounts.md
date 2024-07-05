---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_shared_accounts"
description: |-
  Use this data source to get the list of shared accounts.
---

# huaweicloud_swr_shared_accounts

Use this data source to get the list of shared accounts.

## Example Usage

```hcl
variable "organization" {}
variable "repository" {}

data "huaweicloud_swr_shared_accounts" "test" {
  organization = var.organization
  repository   = var.repository
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `organization` - (Required, String) Specifies the name of the organization to which the repository belongs.

* `repository` - (Required, String) Specifies the name of the repository.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `shared_accounts` - All shared accounts that match the filter parameters.
  The [shared_accounts](#swr_shared_accounts) structure is documented below.

<a name="swr_shared_accounts"></a>
The `repositories` block supports:

* `organization` - The name of the organization to which the repository belongs.

* `repository` - The name of the repository.

* `shared_account` - The shared account name.

* `status` - Whether the sharing account is valid.

* `permit` - The permissions of the shared account.

* `deadline` - The expiration time.

* `description` - The description.

* `creator_id` - The creator ID.

* `created_by` - The name of the creator.

* `created_at` - The creation time.

* `updated_at` - The update time.
