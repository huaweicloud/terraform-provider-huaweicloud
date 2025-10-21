---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_scheduled_events"
description: |-
  Use this data source to get the list of scheduled events of an instance.
---

# huaweicloud_compute_scheduled_events

Use this data source to get the list of scheduled events of an instance.

## Example Usage

```hcl
data "huaweicloud_compute_scheduled_events" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `event_id` - (Optional, String) Specifies the event ID.

* `instance_id` - (Optional, List) Specifies the instance ID.

* `type` - (Optional, List) Specifies the event type.

* `state` - (Optional, List) Specifies the event status.
  Value options:
  + **inquiring**: An event is waiting to be authorized with the start time specified.
    The system will complete operations within a specified time. For details, see the response event.
  + **scheduled**: The event is waiting for the system to schedule resources.
  + **executing**: The system has scheduled resources and is rectifying the fault.
  + **completed**: The system has completed the event execution.
    Check the impacts on services. If any problems occur, contact technical support.
  + **failed**: The system fails to automatically rectify the fault.
  + **canceled**: The event has been canceled by the system.

* `publish_since` - (Optional, String) Specifies the start time of publishing an event. The value is filtered by time range.
  The value is in the format of **yyyy-MM-ddTHH:mm:ssZ** in UTC+0 and complies with ISO8601.

* `publish_until` - (Optional, String) Specifies the end time of publishing an event. The value is filtered by time range.
  The value is in the format of **yyyy-MM-ddTHH:mm:ssZ** in UTC+0 and complies with ISO8601.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `events` - Indicates the event list.

  The [events](#events_struct) structure is documented below.

<a name="events_struct"></a>
The `events` block supports:

* `id` - Indicates the event ID.

* `state` - Indicates the event status.

* `type` - Indicates the event type.

* `description` - Indicates the event description.

* `publish_time` - Indicates when the event is published.

* `start_time` - Indicates the event start time.

* `finish_time` - Indicates the event completion time.

* `instance_id` - Indicates the instance ID.

* `not_before` - Indicates when the event is scheduled to start.

* `not_after` - Indicates when the event is scheduled to end.

* `not_before_deadline` - Indicates the deadline of starting a scheduled event.

* `execute_options` - Indicates the event execution option.

  The [execute_options](#events_execute_options_struct) structure is documented below.

* `source` - Indicates the event source.

  The [source](#events_source_struct) structure is documented below.

<a name="events_execute_options_struct"></a>
The `execute_options` block supports:

* `device` - Indicates the device name of the local disk.

* `wwn` - Indicates the unique ID used for attaching the local disk.

* `serial_number` - Indicates the SN of the local disk.

* `resize_target_flavor_id` - Indicates the flavor ID.

* `migrate_policy` - Indicates the instance migration policy.

* `executor` - Indicates the executor.

<a name="events_source_struct"></a>
The `source` block supports:

* `type` - Indicates the source type of the scheduled event.
  The value can be:
  + **platform**: The event is initiated by the O&M platform.
  + **user**: The event is initiated by users.

* `host_scheduled_event_id` - Indicates the ID of the scheduled event for the host.
