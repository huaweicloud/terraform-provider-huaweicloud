---
subcategory: "Cloud Connect (CC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cc_global_connection_bandwidth_configs"
description: |-
  Use this data source to get the tenant configuration of a global connection bandwidth.
---

# huaweicloud_cc_global_connection_bandwidth_configs

Use this data source to get the tenant configuration of a global connection bandwidth.

## Example Usage

```hcl
data "huaweicloud_cc_global_connection_bandwidth_configs" "test"{}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `configs` - Indicates the dynamic configuration items for purchasing a global connection bandwidth.

  The [configs](#configs_struct) structure is documented below.

<a name="configs_struct"></a>
The `configs` block supports:

* `bind_limit` - Indicates the maximum number of instances that are allowed to use a shared bandwidth.

* `enable_change_95` - Indicates whether standard 95th percentile bandwidth billing can be changed to billing by
  bandwidth capacity.

* `size_range` - Indicates the capacity of global connection bandwidths by billing option.

  The [size_range](#configs_size_range_struct) structure is documented below.

* `services` - Indicates the instance type.

* `gcb_type` - Indicates the bandwidth type.

* `ratio_95peak_plus` - Indicates the percentage of the minimum bandwidth in enhanced 95th percentile bandwidth billing.

* `ratio_95peak_guar` - Indicates the percentage of the minimum bandwidth in standard 95th percentile bandwidth billing.

* `sla_level` - Indicates the line grade.

* `enable_spec_code` - Indicates whether multiple line specifications are supported.

* `charge_mode` - Indicates the list of supported billing options.

* `crossborder` - Indicates whether a cross-border permit is approved.

* `quotas` - Indicates the quota information.

  The [quotas](#configs_quotas_struct) structure is documented below.

* `enable_area_bandwidth` - Indicates whether to enable the geographic region bandwidth.

<a name="configs_size_range_struct"></a>
The `size_range` block supports:

* `type` - Indicates the billing option of a global connection bandwidth.
  The value can be:
  + **bwd**: billing by bandwidth capacity
  + **95**: standard 95th percentile bandwidth billing
  + **95avr**: average daily 95th percentile bandwidth

* `min` - Indicates the minimum global connection bandwidth, in Mbit/s.

* `max` - Indicates the maximum global connection bandwidth, in Mbit/s.

<a name="configs_quotas_struct"></a>
The `quotas` block supports:

* `quota` - Indicates the quotas.

* `used` - Indicates the used quotas.

* `type` - Indicates the quota type.
  The value can be:
  + **gcb.size**: global connection bandwidth capacity
  + **gcb.count**: number of global connection bandwidths
