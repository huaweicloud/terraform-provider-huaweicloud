---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_event_report"
description: |-
  Use this resource to report AOM event or alarm within HuaweiCloud.
---

# huaweicloud_aom_event_report

Use this resource to report AOM event or alarm within HuaweiCloud.

-> This resource is only a one-time action resource for report error event or alarm. Deleting this resource
   will not clear the corresponding request record, but will only remove the resource information from
   the tfstate file.

## Example Usage

### Report an event or alarm

```hcl
variable "starts_at" {}
variable "event_name" {}
variable "event_severity" {}
variable "event_type" {}
variable "resource_id" {}
variable "resource_type" {}
variable "resource_provider" {}
variable "message" {}

resource "huaweicloud_aom_event_report" "test" {
  events {
    starts_at = var.starts_at

    metadata = {
      event_name        = var.event_name
      event_severity    = var.event_severity
      event_type        = var.event_type
      resource_id       = var.resource_id
      resource_type     = var.resource_type
      resource_provider = var.resource_provider
    }

    annotations = jsonencode({
      message = var.message
    })
  }
}
```

### Clear an alarm

```hcl
variable "action" {}
variable "event_name" {}
variable "event_severity" {}
variable "event_type" {}
variable "resource_id" {}
variable "resource_type" {}
variable "resource_provider" {}
variable "ends_at" {}

resource "huaweicloud_aom_event_report" "clear_alarm" {
  action = var.action

  events {
    metadata = {
      event_name        = var.event_name
      event_severity    = var.event_severity
      event_type        = var.event_type
      resource_id       = var.resource_id
      resource_type     = var.resource_type
      resource_provider = var.resource_provider
    }

    ends_at = var.ends_at
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the event source are located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `events` - (Required, List, NonUpdatable) Specifies the list of events or alarms to be reported.  
  The [events](#aom_events) structure is documented below.

* `action` - (Optional, String, NonUpdatable) Specifies the action type of the request.  
  Default value is empty, which means reporting events or alarms.  
  Set to **clear** to clear alarms.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID to which the resource
  belongs.

<a name="aom_events"></a>
The `events` block supports:

* `metadata` - (Required, Map, NonUpdatable) Specifies the detail of the event or alarm, in key:value pair format.
  The following fields are required:
  + `event_name` - (Required, String) Specifies the name of the event or alarm.
  + `event_severity` - (Required, String) Specifies the severity level of the event or alarm.
    The valid values are as follows:
    - **Critical**
    - **Major**
    - **Minor**
    - **Info**
  + `event_type` - (Required, String) Specifies the type of the event or alarm.
    The valid values are as follows:
    - **event**
    - **alarm**
  + `resource_provider` - (Required, String) Specifies the cloud service name corresponding to the event.
  + `resource_type` - (Required, String) Specifies the resource type corresponding to the event.
  + `resource_id` - (Required, String) Specifies the resource information corresponding to the event.

  The value length of each metadata field is `1` to `2,048` characters.

* `starts_at` - (Optional, Int, NonUpdatable) Specifies the time when the event or alarm occurred, in UTC
  milliseconds timestamp.  
  For example: 2024-10-16 16:03:01 needs to be converted to UTC milliseconds timestamp: `1702759381000`.  
  Required if the `action` is empty.

* `ends_at` - (Optional, Int, NonUpdatable) Specifies the time when the event or alarm was cleared, in UTC
  milliseconds timestamp.  
  Default value is `0`, indicating that the alarm was not cleared.  
  Required if the `action` is **clear**.

* `timeout` - (Optional, Int, NonUpdatable) Specifies the automatic clearing time for expired alarms, in milliseconds.
  The maximum clearing time is `15` days.  
  + When `action` is empty:
    - If `timeout` is empty, the default clearing time is `15` days.
    - If `timeout` is specified with a time, the format of time should format be the corresponding milliseconds.  
      For example, if you want to set the clearing time to `5` days, the corresponding milliseconds: `432000000`.
  + When `action` is **clear**:
    - This parameter is not required when clearing alarms.

* `annotations` - (Optional, String, NonUpdatable) Specifies the additional fields of the event or alarm, in JSON
  format.  
  This parameter can be empty.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
