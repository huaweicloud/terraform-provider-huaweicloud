---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_catalog_statistics"
description: |-
  Use this data source to get the catalog statistics within HuaweiCloud.
---

# huaweicloud_dsc_catalog_statistics

Use this data source to get the catalog statistics within HuaweiCloud.

## Example Usage

```hcl
variable "type_id" {}

data "huaweicloud_dsc_catalog_statistics" "test" {
  type_id = var.type_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `label_id` - (Optional, String) Specifies the group label ID used to filter the statistics of a specific group.  
  The `label_id` and `type_id` must be specified at least one.

* `type_id` - (Optional, String) Specifies the type ID used to filter the statistics of a specific type.  
  The `label_id` and `type_id` must be specified at least one.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `bucket` - The bucket risk statistics information.

  The [bucket](#obs_static_info_struct) structure is documented below.

* `column` - The column risk statistics information.

  The [column](#column_static_info_struct) structure is documented below.

* `database` - The database risk statistics information.

  The [database](#static_info_struct) structure is documented below.

* `file` - The file risk statistics information.

  The [file](#obs_static_info_struct) structure is documented below.

* `table` - The table risk statistics information.

  The [table](#static_info_struct) structure is documented below.

<a name="obs_static_info_struct"></a>
The `bucket` and `file` blocks support:

* `risk_num` - The risk number.

* `total` - The total number.

<a name="column_static_info_struct"></a>
The `column` block supports:

* `risk_num` - The risk number.

* `total` - The total number.

* `week_on_week` - The week-on-week ratio.

<a name="static_info_struct"></a>
The `database` and `table` blocks support:

* `risk_num` - The risk number.

* `total` - The total number.

* `week_on_week` - The week-on-week ratio.
