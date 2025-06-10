---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_product_infos"
description: |-
  Use this data source to get the list of HSS product information within HuaweiCloud.
---

# huaweicloud_hss_product_infos

Use this data source to get the list of HSS product information within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_product_infos" "test" {
  site_code = "HWC_CN"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `site_code` - (Optional, String) Specifies the site information. The valid values are as follows:
  + **HWC_CN**: Chinese mainland.
  + **HWC_HK**: International.
  + **HWC_EU**: Europe.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the resource belongs.
  This parameter is valid only when the enterprise project function is enabled.
  For enterprise users, if omitted, default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The product information list.
  The [data_list](#product_info_structure) structure is documented below.

<a name="product_info_structure"></a>
The `data_list` block supports:

* `charging_mode` - The billing modes. The valid values are as follows:
  + **packet_cycle**: Yearly/Monthly subscription.
  + **on_demand**: Pay-per-use.

* `is_auto_renew` - Whether to enable automatic renewal.

* `version_info` - The edition information list.
  The [version_info](#version_info_structure) structure is documented below.

<a name="version_info_structure"></a>
The `version_info` block supports:

* `version` - The HSS edition. The value can be:
  + **hss.version.basic**: Basic edition.
  + **hss.version.advanced**: Professional edition.
  + **hss.version.enterprise**: Enterprise edition.
  + **hss.version.premium**: Premium edition.
  + **hss.version.wtp**: WTP edition.
  + **hss.version.container.enterprise**: Container edition.
  + **hss.version.small.hsp**: Small HSP edition.
  + **hss.version.imagescan**: Image scan edition.
  + **hss.version.antivirus**: Antivirus edition.

* `periods` - The period information list.
  The [periods](#period_info_structure) structure is documented below.

<a name="period_info_structure"></a>
The `periods` block supports:

* `period_vals` - Value string of the required duration. Multiple values are separated by commas (,).
  For example: "1,2,3,4,5,6,7,8,9" for monthly subscription or "1,2,3,5" for yearly subscription.

* `period_unit` - Required duration unit. The valid values are as follows:
  + **year**: Year.
  + **month**: Month.
