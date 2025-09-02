---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_table_histograms"
description: |-
  Use this data source to get the histogram data of a specified table from HuaweiCloud SecMaster.
---

# huaweicloud_secmaster_table_histograms

Use this data source to get the histogram data of a specified table from HuaweiCloud SecMaster.

## Example Usage

```hcl
variable "workspace_id" {}
variable "table_id" {}

data "huaweicloud_secmaster_table_histograms" "test" {
  workspace_id = var.workspace_id
  table_id     = var.table_id
  query        = "*"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `table_id` - (Required, String) Specifies the ID of the table to get histogram data.

* `query` - (Required, String) Specifies the search query string.

* `from` - (Optional, Int) Specifies the start timestamp in milliseconds.

* `to` - (Optional, Int) Specifies the end timestamp in milliseconds.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `histograms` - The list of histogram data.
  The [histogram](#histogram_struct) structure is documented below.

<a name="histogram_struct"></a>
The `histogram` block supports:

* `count` - The count of logs in the time range.

* `from` - The start timestamp of the time range.

* `to` - The end timestamp of the time range.
