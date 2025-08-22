---
subcategory: "Event Grid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_event_batch_action"
description: |-
  Use this resource to publish events to an event channel within HuaweiCloud.
---

# huaweicloud_eg_event_batch_action

Use this resource to publish events to an event channel within HuaweiCloud.

-> This resource is only a one-time action resource for dispatch events to channel. Deleting this resource will not
   clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "event_channel_id" {}
variable "events" {
  type = list(object({
    id                = string
    source            = string
    spec_version      = string
    type              = string
    data_content_type = string
    data              = string
    time              = string
    subject           = string
  }))
}

resource "huaweicloud_eg_event_batch_action" "test" {
  channel_id = var.event_channel_id

  dynamic "events" {
    for_each = var.events

    content {
      id                = events.value.id
      source            = events.value.source
      spec_version      = events.value.spec_version
      type              = events.value.type
      data_content_type = events.value.data_content_type
      data              = events.value.data
      time              = events.value.time
      subject           = events.value.subject
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the event channel is located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `channel_id` - (Required, String, NonUpdatable) Specifies the ID of the event channel.

* `events` - (Required, List, NonUpdatable) Specifies the list of events to be published.  
  The [events](#eg_events_attr) structure is documented below.

<a name="eg_events_attr"></a>
The `events` block supports:

* `id` - (Required, String) Specifies the ID of the event.

* `source` - (Required, String) Specifies the name of the event source.  
  For the detail, please following [reference documentation](https://tools.ietf.org/html/rfc3986#section-4.1)

* `spec_version` - (Required, String) Specifies the CloudEvents protocol version.  
  The spec version must follow the pattern `major.minor`

* `type` - (Required, String) Specifies the type of the event.

* `data_content_type` - (Optional, String) Specifies the content type of the event data.  
  For the detail, please following [reference documentation](https://tools.ietf.org/html/rfc2046)

* `data_schema` - (Optional, String) Specifies the URI of the event data schema.  
  For the detail, please following [reference documentation](https://tools.ietf.org/html/rfc3986#section-4.3)

* `data` - (Optional, String) Specifies the payload content of the event, in JSON format.  
  The content of data must follow the data schema description.

* `time` - (Optional, String) Specifies the time when the event occurred, in UTC format.

* `subject` - (Optional, String) Specifies the subject of the event.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
