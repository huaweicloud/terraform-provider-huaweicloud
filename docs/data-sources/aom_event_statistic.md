---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_event_statistic"
description: |-
  Use this data source to get the event and alarm statistics.
---

# huaweicloud_aom_event_statistic

Use this data source to get the event and alarm statistics.

## Example Usage

```hcl
data "huaweicloud_aom_event_statistic" "test" {
  time_range = "-1.-1.60"
  step       = 60000
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `time_range` - (Required, String) Specifies the time range for querying event and alarm statistics.  
  The format is **startTimeInMillis.endTimeInMillis.durationInMinutes**, for example:
  + **-1.-1.60**: Query the last 60 minutes
  + **1650852000000.1650852300000.5**: Query from 2022-04-25 10:00:00 to 2022-04-25 10:05:00

* `step` - (Optional, Int) Specifies the statistical step size in milliseconds.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `step_result` - The statistical step size in milliseconds.

* `timestamps` - The time series corresponding to the statistical results.

* `series` - The statistical results for different severity levels at the same time series.  
  The [series](#aom_event_statistic_series) structure is documented below.

* `summary` - The summary of various alarm information quantities.

<a name="aom_event_statistic_series"></a>
The `series` block supports:

* `event_severity` - The event or alarm severity level.

* `values` - The statistical results for events or alarms at each time point.
