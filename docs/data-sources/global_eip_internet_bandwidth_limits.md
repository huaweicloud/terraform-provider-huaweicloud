---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eip_internet_bandwidth_limits"
description: |-
  Use this data source to get the list of global EIP internet bandwidth limits.
---

# huaweicloud_global_eip_internet_bandwidth_limits

Use this data source to get the list of global EIP internet bandwidth limits.

## Example Usage

```hcl
data "huaweicloud_global_eip_internet_bandwidth_limits" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fields` - (Optional, List) Specifies the fields to return.  
  Supported values include **id**, **charge_mode**, **min_size**, **ext_limit**, **max_size**, and **type**.

* `sort_key` - (Optional, String) Specifies the sort fields.  
  Valid value is **id**.

* `sort_dir` - (Optional, String) Specifies the sort directions.  
  Valid values are **asc** and **desc**.

* `charge_mode` - (Optional, String) Specifies the billing modes used to filter the results.  
  Valid values are **bandwidth**, **95peak_bidirection**, **95peak_plus_1000** and **95peak_guar**.

* `type` - (Optional, String) Specifies the global EIP internet bandwidth type used to filter the results.
  The valid values are as follows:
  + **BioDir**
  + **Enter**
  + **Flx**
  + **standard**
  + **Ext**
  + **Ext_Enter**
  + **Ext_BioDir**
  + **IPv6**
  + **Zixun**
  + **STD**
  + **Lite**
  + **Spec**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `internet_bandwidth_limits` - The list of global EIP internet bandwidth limits.

  The [internet_bandwidth_limits](#internet_bandwidth_limits_struct) structure is documented below.

<a name="internet_bandwidth_limits_struct"></a>
The `internet_bandwidth_limits` block supports:

* `id` - The global EIP internet bandwidth limit ID.

* `charge_mode` - The billing mode.

* `min_size` - The minimum purchasable bandwidth size for this global EIP internet bandwidth type.

* `ext_limit` - The auxiliary limit information.

  The [ext_limit](#ext_limit_struct) structure is documented below.

* `max_size` - The maximum purchasable bandwidth size for this global EIP internet bandwidth type.

* `type` - The global EIP internet bandwidth type.

<a name="ext_limit_struct"></a>
The `ext_limit` block supports:

* `min_ingress_size` - The minimum ingress global EIP internet bandwidth size.

* `max_ingress_size` - The maximum ingress global EIP internet bandwidth size.

* `ratio_95peak` - The enhanced 95th percentile EIP guaranteed ratio.
