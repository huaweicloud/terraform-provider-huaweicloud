---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_alarm_handling_trend"
description: |-
  Use this data source to get the handling statistics of alarms and events within HuaweiCloud.
---

# huaweicloud_dsc_alarm_handling_trend

Use this data source to get the handling statistics of alarms and events within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_alarm_handling_trend" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `time_axis` - The time axis.

* `total_alarm_variation` - The variation list of the total number of alarms.

* `total_event_variation` - The variation list of the total number of events.

* `untreated_alarm_variation` - The variation list of the number of untreated alarms.

* `untreated_event_variation` - The variation list of the number of untreated events.
