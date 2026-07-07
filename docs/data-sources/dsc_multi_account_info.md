---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_multi_account_info"
description: |-
  Use this data source to get the organization information of the current account within HuaweiCloud.
---

# huaweicloud_dsc_multi_account_info

Use this data source to get the organization information of the current account within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_multi_account_info" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `is_admin` - Whether the account is an administrator.

* `is_delegated_admin` - Whether the account is a delegated administrator.

* `is_trusted_service` - Whether the account is a trusted service.
