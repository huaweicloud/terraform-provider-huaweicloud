---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_supported_column_types"
description: |-
  Use this data source to get the supported column types within HuaweiCloud.
---

# huaweicloud_dsc_supported_column_types

Use this data source to get the supported column types within HuaweiCloud.

## Example Usage

```hcl
variable "data_source_type" {}

data "huaweicloud_dsc_supported_column_types" "test" {
  data_source_type = var.data_source_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `data_source_type` - (Required, String) Specifies the data source type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_source_type_response` - The data source type in the response.

* `supported_column_types` - The supported column type list.
