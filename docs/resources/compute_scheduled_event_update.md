---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_scheduled_event_update"
description: |-
  Manages an ECS scheduled event execution start time update resource within HuaweiCloud.
---

# huaweicloud_compute_scheduled_event_update

Manages an ECS scheduled event execution start time update resource within HuaweiCloud.

## Example Usage

```hcl
variable "event_id" {}

resource "huaweicloud_compute_scheduled_event_update" "test" {
  event_id   = var.event_id
  not_before = "2025-07-09T10:40:00Z"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `event_id` - (Required, String, NonUpdatable) Specifies the event ID.

* `not_before` - (Required, String, NonUpdatable) Specifies the scheduled start time. The new start time must be earlier
  than the deadline of the scheduled event start time

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the event ID.
