---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_multi_organizations"
description: |-
  Use this data source to get the list of organizations within HuaweiCloud.
---

# huaweicloud_dsc_multi_organizations

Use this data source to get the list of organizations within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_multi_organizations" "test" {}
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

* `ou_list` - The organizational unit list.

  The [ou_list](#ou_list_struct) structure is documented below.

<a name="ou_list_struct"></a>
The `ou_list` block supports:

* `id` - The organizational unit ID.

* `name` - The organizational unit name.
