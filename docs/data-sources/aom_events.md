---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_events"
description: |-
  Use this data source to get the list of events and alarms.
---

# huaweicloud_aom_events

Use this data source to get the list of events and alarms.

## Example Usage

```hcl
data "huaweicloud_aom_events" "test" {
  time_range = "-1.-1.60"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `time_range` - (Required, String) Specifies the time range for querying events and alarms.  
  The format is **startTimeInMillis.endTimeInMillis.durationInMinutes**, for example:
  + **-1.-1.60**: Query the last 60 minutes
  + **1650852000000.1650852300000.5**: Query from `2022-04-25 10:00:00` to `2022-04-25 10:05:00`

* `step` - (Optional, Int) Specifies the statistical step size in milliseconds.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `events` - The list of events and alarms that matched the filter parameters.  
  The [events](#aom_events) structure is documented below.

<a name="aom_events"></a>
The `events` block supports:

* `id` - The ID of the event or alarm.

* `event_sn` - The alarm serial number.

* `starts_at` - The time when the event or alarm occurred, CST millisecond timestamp.

* `ends_at` - The time when the event or alarm was cleared, CST millisecond timestamp, **0** means not cleared.

* `arrives_at` - The time when the event arrived at the system, CST millisecond timestamp.

* `timeout` - The automatic clearing time for alarms in milliseconds.

* `enterprise_project_id` - The enterprise project ID to which the event or alarm belongs.

* `metadata` - The detailed information (key/value pair) of the event or alarm.

* `annotations` - The additional fields of the event or alarm, in JSON format.

* `policy` - The open alarm policy, in JSON format.
