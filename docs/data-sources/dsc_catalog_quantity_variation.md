---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_catalog_quantity_variation"
description: |-
  Use this data source to get the catalog quantity variation within HuaweiCloud.
---

# huaweicloud_dsc_catalog_quantity_variation

Use this data source to get the catalog quantity variation within HuaweiCloud.

## Example Usage

```hcl
variable "type_id" {}

data "huaweicloud_dsc_catalog_quantity_variation" "test" {
  type_id = var.type_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `label_id` - (Optional, String) Specifies the label ID used to filter the data quantity variation.
  Either `label_id` or `type_id` must be specified.

* `type_id` - (Optional, String) Specifies the type ID used to filter the data quantity variation.
  Either `label_id` or `type_id` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `sensitive_number_variation` - The sensitive number variation trend.

* `time_axis` - The time axis.

* `total_number_variation` - The total number variation trend.
