---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_overviews_bandwidth_timeline"
description: |-
  Use this data source to query the average bandwidth usage.
---

# huaweicloud_waf_overviews_bandwidth_timeline

Use this data source to query the average bandwidth usage.

## Example Usage

```hcl
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_waf_overviews_bandwidth_timeline" "test" {
  from = var.start_time
  to   = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `from` - (Required, Int) Specifies the query start time.
  The format is 13-digit timestamp in millisecond.

* `to` - (Required, Int) Specifies the query end time.
  The format is 13-digit timestamp in millisecond.

* `group_by` - (Optional, String) Specifies the display dimension.
  The value can be **DAY**, indicates data is displayed by the day.
  If this parameter is not specified, the data is displayed by the minute.

* `display_option` - (Optional, String) Specifies the number of sent or received bytes.
  The valid values are as follows:
  + **1**: Indicates view the peak value.
  + **0**: Indicates view the average value.

* `hosts` - (Optional, String) Specifies the ID of the domain.

* `instances` - (Optional, String) Specifies the ID of the dedicated WAF instance.
  This parameter is used to query the average bandwidth usage of domain protected by the dedicated WAF instance.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  If you want to query resources under all enterprise projects, set this parameter to **all_granted_eps**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `bandwidths` - The bandwidth statistics over the time.

  The [bandwidths](#bandwidths_struct) structure is documented below.

<a name="bandwidths_struct"></a>
The `bandwidths` block supports:

* `key` - The key type.
  The options are **BANDWIDTH**, **IN_BANDWIDTH** and **OUT_BANDWIDTH**.

* `timeline` - The statistics data over time for the corresponding key.

  The [timeline](#bandwidths_timeline_struct) structure is documented below.

<a name="bandwidths_timeline_struct"></a>
The `timeline` block supports:

* `time` - The time point.

* `num` - The statistics data for the time range from the previous time point to the point specified by `time`.
