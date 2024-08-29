---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_events"
description: |-
  Use this data source to get the list of CES events.
---

# huaweicloud_ces_events

Use this data source to get the list of CES events.

## Example Usage

```hcl
data "huaweicloud_ces_events" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the event type.
  The value can be **EVENT.SYS** (system event) or **EVENT.CUSTOM** (custom event).

* `name` - (Optional, String) Specifies the event name.

* `from` - (Optional, String) Specifies the start time of the query.
  The time is in UTC. The format is **yyyy-MM-dd HH:mm:ss**.
  The start time cannot be greater than the current time.

* `to` - (Optional, String) Specifies the end time of the query.
  The time is in UTC. The format is **yyyy-MM-dd HH:mm:ss**.
  The start time needs to be less than the end time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `events` - The event records.

  The [events](#events_struct) structure is documented below.

<a name="events_struct"></a>
The `events` block supports:

* `event_name` - The event name.

* `event_type` - The event type.

* `event_count` - The number of occurrences of this event within the specified query time range.

* `latest_occur_time` - The time when the event last occurred. The time is in UTC.

* `latest_event_source` - The latest event source.
