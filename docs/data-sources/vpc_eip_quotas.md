---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_eip_quotas"
description: |-
  Use this data source to get the list of EIP related quotas within HuaweiCloud.
---

# huaweicloud_vpc_eip_quotas

Use this data source to get the list of EIP related quotas within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_vpc_eip_quotas" "test" {}
```
## Argument Reference

The following arguments are supported:

* `region`- (Optional, String) Specifies the region in which to query the resource.
If omitted, the provider-level region will be used.

* `tags` - (Optional, String) Specifies the quota type to filter.
The valid values are publicip, bandwidth.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas`- The list of EIP quotas.

The [quotas](#quotas_struct) structure is documented below.

<a name="quotas_struct"></a>
The `quotas` block supports:

* `type` - The type of the EIP quota.

* `used` - The number of used EIP quotas.

* `quota` - The total number of EIP quotas.value of the quota that can be modified.