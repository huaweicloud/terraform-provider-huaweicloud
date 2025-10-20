---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_scheduled_event_accept"
description: |-
  Manages an ECS respond to the system maintenance events in pending authorization state resource within HuaweiCloud.
---

# huaweicloud_compute_scheduled_event_accept

Manages an ECS respond to the system maintenance events in pending authorization state resource within HuaweiCloud.

## Example Usage

```hcl
variable "event_id" {}

resource "huaweicloud_compute_scheduled_event_accept" "test" {
  event_id = var.event_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `event_id` - (Required, String, NonUpdatable) Specifies the event ID.

* `not_before` - (Optional, String, NonUpdatable) Specifies the scheduled start time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the event ID.
