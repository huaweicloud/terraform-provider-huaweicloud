---
layout: "huaweicloud"
page_title: "Flexibleengine: huaweicloud_dcs_product_v1"
sidebar_current: "docs-huaweicloud-datasource-dcs-product-v1"
description: |-
  Get information on an Flexibleengine dcs product.
---

# huaweicloud\_dcs\_product

Use this data source to get the ID of an available Flexibleengine dcs product.
This is an alternative to `huaweicloud_dcs_product_v1`

## Example Usage

```hcl

data "huaweicloud_dcs_product" "product1" {
  spec_code = "dcs.single_node"
}
```

## Argument Reference

* `spec_code` - (Optional) DCS instance specification code. For details, see
[Querying Service Specifications](https://support.huaweicloud.com/en-us/api-dcs/dcs-api-0312040.html).


## Attributes Reference

`id` is set to the ID of the found product. In addition, the following attributes
are exported:

* `spec_code` - See Argument Reference above.
