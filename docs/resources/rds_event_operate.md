---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_event_operate"
description: |-
  Manages an RDS event operate resource within HuaweiCloud.
---

# huaweicloud_rds_event_operate

Manages an RDS event operate resource within HuaweiCloud.

## Example Usage

```hcl
variable "event_id" {}
variable "instance_id" {}

resource "huaweicloud_rds_event_operate" "test" {
  event_instances {
    event_id    = var.event_id
    instance_id = var.instance_id
  }

  operation_type = "cancel"

  event_schedule_window {
    planned_day = "2025-11-13"
    start_time  = "06:00"
    end_time    = "10:00"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `event_instances` - (Required, List, NonUpdatable) Specifies the list of event instances to operate.
  The [event_instances](#event_instances_struct) structure is documented below.

* `operation_type` - (Required, String, NonUpdatable) Specifies the event operation type.
  The valid values are as follows:
  + **cancel**: Cancel the schedule event.
  + **execute**: Execute the schedule event.
  + **reservation**: Reserve the schedule event.

* `event_schedule_window` - (Optional, List, NonUpdatable) Specifies the event schedule window.
  This parameter is required when `operation_type` is set to **reservation**.
  The [event_schedule_window](#event_schedule_window_struct) structure is documented below.

<a name="event_instances_struct"></a>
The `event_instances` block supports:

* `event_id` - (Required, String, NonUpdatable) Specifies the event ID.

* `instance_id` - (Required, String, NonUpdatable) Specifies the instance ID.

<a name="event_schedule_window_struct"></a>
The `event_schedule_window` block supports:

* `planned_day` - (Required, String, NonUpdatable) Specifies the planned execution date.
  The format is **yyyy-MM-dd**.

* `start_time` - (Optional, String, NonUpdatable) Specifies the planned execution window start time.
  The format is **hh:mm**.

* `end_time` - (Optional, String, NonUpdatable) Specifies the planned execution window end time.
  The format is **hh:mm**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the event ID.

* `results` - The list of event operation results.
  The [results](#schedule_event_operation_results) structure is documented below.

<a name="schedule_event_operation_results"></a>
The `results` block supports:

* `id` - The event ID.

* `instance_id` - The instance ID.

* `job_id` - The job ID.

* `error_code` - The error code.

* `error_msg` - The error message.

* `success` - Whether the operation is successful.
