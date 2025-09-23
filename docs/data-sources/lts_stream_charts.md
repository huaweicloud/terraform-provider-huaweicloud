---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_stream_charts"
description: |-
  Use this data source to get the list of LTS stream charts.
---

# huaweicloud_lts_stream_charts

Use this data source to get the list of LTS stream charts.

## Example Usage

```hcl
variable "group_id" {}
variable "stream_id" {}

data "huaweicloud_lts_stream_charts" "test" {
  log_group_id  = var.group_id
  log_stream_id = var.stream_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the chart list.  
  If omitted, the provider-level region will be used.

* `log_group_id` - (Required, String) Specifies the ID of the log group.

* `log_stream_id` - (Required, String) Specifies the ID of the log stream.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `charts` - The list of charts.  
  The [charts](#stream_charts_attr) structure is documented below.

<a name="stream_charts_attr"></a>
The `charts` block supports:

* `id` - The ID of the chart.

* `name` - The name of the chart.

* `type` - The type of the chart.

* `log_group_id` - The ID of the log group.

* `log_group_name` - The name of the log group.

* `log_stream_id` - The ID of the log stream.

* `log_stream_name` - The name of the log stream.

* `sql` - The SQL statement of the chart.

* `config` - The configuration of the chart.  
  The [config](#stream_charts_config_attr) structure is documented below.

<a name="stream_charts_config_attr"></a>
The `config` block supports:

* `page_size` - The page size of the chart.

* `can_sort` - Whether to enable sorting.

* `can_search` - Whether to enable search.
