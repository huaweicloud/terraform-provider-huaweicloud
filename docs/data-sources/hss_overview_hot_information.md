---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_overview_hot_information"
description: |-
  Use this data source to get the list of HSS overview hot information within HuaweiCloud.
---

# huaweicloud_hss_overview_hot_information

Use this data source to get the list of HSS overview hot information within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_overview_hot_information" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number.

* `data_list` - The list of hot information.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `title` - The title of the hot information.

* `update_time` - The update time of the hot information.

* `severity_level` - The severity level of the hot information.
