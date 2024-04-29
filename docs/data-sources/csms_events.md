---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_events"
description: ""
---

# huaweicloud_csms_events

Use this data source to get the list of CSMS events.

## Example Usage

```hcl
variable "event_name" {}

data "huaweicloud_csms_events" "test" {
  name = var.event_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the event.

* `event_id` - (Optional, String) Specifies the ID of the event.

* `status` - (Optional, String) Specifies the event status. Valid values are **ENABLED** and **DISABLED**.
  Only the event in **ENABLED** status can be triggered.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `events` - Indicates the events list.
  The [events](#CSMS_events) structure is documented below.

<a name="CSMS_events"></a>
The `events` block supports:

* `name` - Indicates the event name.

* `event_id` - The event ID.

* `event_types` - Indicates the event type. Valid values are:
  + **SECRET_VERSION_CREATED**: Triggered when a version of a secret is created.
  + **SECRET_VERSION_EXPIRED**: Triggered when a secret version expires, and only once per expiration.
  + **SECRET_ROTATED**: Triggered when a secret is rotated. Currently, only RDS secrets can be automatically rotated.
  + **SECRET_DELETED**: Triggered when a secret is deleted.

* `status` - Indicates the event status.

* `created_at` - Indicates the time when the event created, in UTC format.

* `updated_at` - Indicates the time when the event updated, in UTC format.

* `notification` - Indicates the event notification list.
  The [notification](#notification) structure is documented below.

<a name="notification"></a>
The `notification` block supports:

* `target_type` - Indicates the object type of the event notification.
* `target_id` - Indicates the object ID of the event notification.
* `target_name` - Indicates the object name of the event notification.
