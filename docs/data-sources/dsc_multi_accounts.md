---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_multi_accounts"
description: |-
  Use this data source to get the list of organization accounts within HuaweiCloud.
---

# huaweicloud_dsc_multi_accounts

Use this data source to get the list of organization accounts within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_multi_accounts" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `accounts` - The account information list.

  The [accounts](#accounts_struct) structure is documented below.

<a name="accounts_struct"></a>
The `accounts` block supports:

* `domain_id` - The domain ID.

* `name` - The account name.

* `project_id` - The project ID.
