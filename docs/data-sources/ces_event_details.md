---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_event_details"
description: |-
  Use this data source to get the CES event details.
---

# huaweicloud_ces_event_details

Use this data source to get the CES event details.

## Example Usage

```hcl
variable "name" {}
variable "type" {}

data "huaweicloud_ces_event_details" "test" {
  name = var.name
  type = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Required, String) Specifies the event name.

* `type` - (Required, String) Specifies the event type.
  The value can be **EVENT.SYS** (system event) or **EVENT.CUSTOM** (custom event).

* `source` - (Optional, String) Specifies the event source.

* `level` - (Optional, String) Specifies the event severity.
  The value can be **Critical**, **Major**, **Minor**, or **Info**.

* `user` - (Optional, String) Specifies the name of the user for reporting event monitoring data.

* `state` - (Optional, String) Specifies the event status.
  The value can be **normal**, **warning**, or **incident**.

* `from` - (Optional, String) Specifies the start time of the query.
  The time is in UTC. The format is **yyyy-MM-dd HH:mm:ss**.
  The start time cannot be greater than the current time.

* `to` - (Optional, String) Specifies the end time of the query.
  The time is in UTC. The format is **yyyy-MM-dd HH:mm:ss**.
  The start time needs to be less than the end time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `event_info` - The event information.

  The [event_info](#event_info_struct) structure is documented below.

<a name="event_info_struct"></a>
The `event_info` block supports:

* `event_id` - The event ID.

* `event_name` - The event name.

* `event_source` - The event source.

* `time` - The time when the event occurred.

* `detail` - The event detail.

  The [detail](#event_info_detail_struct) structure is documented below.

<a name="event_info_detail_struct"></a>
The `detail` block supports:

* `event_state` - The event status.

* `event_level` - The event level.

* `event_user` - The event user.

* `content` - The event content.

* `group_id` - The group that the event belongs to.

* `resource_id` - The resource ID.

* `resource_name` - The resource name.

* `event_type` - The event type.

* `dimensions` - The resource dimensions.

  The [dimensions](#detail_dimensions_struct) structure is documented below.

<a name="detail_dimensions_struct"></a>
The `dimensions` block supports:

* `name` - The resource dimension name.

* `value` - The resource dimension value.
