---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_best_practice_account_info"
description: |-
  Use this data source to get the account information used best-practice in Resource Governance Center.
---

# huaweicloud_rgc_best_practice_account_info

Use this data source to get the account information used best-practice in Resource Governance Center.

## Example Usage

```hcl
data "huaweicloud_rgc_best_practice_account_info" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `account_type` - the account type used best-practice.

* `effective_start_time` - The subscription time for best-practice.

* `effective_expiration_time` - The expiration time for best-practice.

* `current_time` - The current time.
