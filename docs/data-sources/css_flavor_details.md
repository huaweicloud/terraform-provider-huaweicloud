---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_flavor_details"
description: |-
  Use this data source to query the flavor dtails.
---

# huaweicloud_css_flavor_details

Use this data source to query the flavor dtails.

## Example Usage

```hcl
variable "flavor_id" {}

data "huaweicloud_css_flavor_details" "test" {
  flavor_id = var.flavor_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `flavor_id` - (Required, String) Specifies the flavor ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `str_id` - The flavor ID.

* `name` - The flavor name.

* `cond_operation_status` - The flavor sales status.
  This parameter takes effect region-wide. If an AZ is not configured in the `cond_operation_az` parameter,
  the value of this parameter is used by default.
  + **normal**: The flavor is in normal commercial use.
  + **sellout**: The flavor has been sold out.

* `cond_operation_az` - The flavor sales status.
  This parameter takes effect AZ-wide

* `flavor_detail` - The flavor attributes.
  The [flavor_detail](#flavor_detail_struct) structure is documented below.

<a name="flavor_detail_struct"></a>
The `flavor_detail` block supports:

* `name` - The attribute name.
  + **cpu**: Flavor CPU attribute.
  + **mem**: Flavor memory attribute.
  + **diskrange**: Flavor disk capacity range attribute.
  + **flavorTypeCn**: Flavor classification attribute in Chinese.
  + **flavorTypeEn**: Flavor classification attribute in English.

* `value` - The attribute value.
