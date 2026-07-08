---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_vpc_ep_quotas"
description: |-
  Use this data source to get the VPC endpoint quotas within HuaweiCloud.
---

# huaweicloud_dsc_vpc_ep_quotas

Use this data source to get the VPC endpoint quotas within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_vpc_ep_quotas" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `quotas` - The VPC endpoint quota information list.
  The [quotas](#dsc_quotas) structure is documented below.

<a name="dsc_quotas"></a>
The `quotas` block supports:

* `quota` - The quota.
* `type` - The quota type.
* `used` - The used quota.
