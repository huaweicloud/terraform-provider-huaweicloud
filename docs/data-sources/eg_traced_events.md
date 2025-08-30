---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_traced_events"
description: |-
  Use this data source to get the list of EG traced events within HuaweiCloud.
---

# huaweicloud_eg_traced_events

Use this data source to get the list of EG traced events within HuaweiCloud.

## Example Usage

### Query all traced events

```hcl
variable "channel_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_eg_traced_events" "test" {
  channel_id = var.channel_id
  start_time = var.start_time
  end_time   = var.end_time
}
```

### Query traced events by event ID

```hcl
variable "channel_id" {}
variable "start_time" {}
variable "end_time" {}
variable "event_id" {}

data "huaweicloud_eg_traced_events" "test" {
  channel_id = var.channel_id
  start_time = var.start_time
  end_time   = var.end_time
  event_id   = var.event_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the events are located.  
  If omitted, the provider-level region will be used.

* `channel_id` - (Required, String) Specifies the ID of the event channel.

* `start_time` - (Required, String) Specifies the start time of the search time range, in UTC format.

* `end_time` - (Required, String) Specifies the end time of the search time range, in UTC format.

* `event_id` - (Optional, String) Specifies the ID of the event.

* `source_name` - (Optional, String) Specifies the name of the event source.

* `event_type` - (Optional, String) Specifies the type of the event.

* `subscription_name` - (Optional, String) Specifies the name of the event subscription.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `events` - The list of traced events that matched filter parameters.  
  The [events](#eg_traced_events_attr) structure is documented below.

<a name="eg_traced_events_attr"></a>
The `events` block supports:

* `id` - The ID of the event.

* `type` - The type of the event.

* `source_name` - The name of the event source.

* `subscription_name` - The name of the event subscription.

* `received_time` - The time when the event to be received, in UTC format.
