---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_multi_account_assets"
description: |-
  Use this data source to get the list of account asset statistics within HuaweiCloud organization.
---

# huaweicloud_dsc_multi_account_assets

Use this data source to get the list of account asset statistics within HuaweiCloud organization.

## Example Usage

```hcl
data "huaweicloud_dsc_multi_account_assets" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `account_id` - (Optional, String) Specifies the account ID used to filter a specific account.

* `parent_id` - (Optional, String) Specifies the parent account ID used to filter child accounts under a specific
  parent account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `accounts` - The account with assets information list.

  The [accounts](#accounts_struct) structure is documented below.

<a name="accounts_struct"></a>
The `accounts` block supports:

* `bigdata_count` - The big data asset count.

* `db_count` - The database asset count.

* `domain_id` - The domain ID.

* `name` - The account name.

* `obs_count` - The object storage asset count.

* `project_id` - The project ID.
