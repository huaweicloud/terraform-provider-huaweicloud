---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_bandwidth_limits"
description: |-
  Use this data source to get a list of bandwidth limits.
---

# huaweicloud_vpc_bandwidth_limits

Use this data source to get a list of bandwidth limits.

## Example Usage

```hcl
data "huaweicloud_vpc_bandwidth_limits" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `charge_mode` - (Optional, String) Specifies the bandwidth charge mode.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `eip_bandwidth_limits` - Indicates the bandwidth limit list.

  The [eip_bandwidth_limits](#eip_bandwidth_limits_struct) structure is documented below.

<a name="eip_bandwidth_limits_struct"></a>
The `eip_bandwidth_limits` block supports:

* `id` - Indicates the bandwidth type ID.

* `charge_mode` - Indicates the bandwidth charging mode.

* `min_size` - Indicates the minimum size that can be purchased for this type of bandwidth.

* `max_size` - Indicates the maximum size that can be purchased for this type of bandwidth.

* `ext_limit` - Indicates the additional restriction information.

  The [ext_limit](#eip_bandwidth_limits_ext_limit_struct) structure is documented below.

<a name="eip_bandwidth_limits_ext_limit_struct"></a>
The `ext_limit` block supports:

* `min_ingress_size` - Indicates the minimum cloud access rate limit.

* `max_ingress_size` - Indicates the maximum cloud access rate limit.

* `ratio_95peak` - Indicates the 95 Minimum charging rate.
