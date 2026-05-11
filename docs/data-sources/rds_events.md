---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_schedule_events"
description: |-
  Use this data source to query the schedule events within HuaweiCloud.
---

# huaweicloud_rds_schedule_events

Use this data source to query the schedule events within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_rds_schedule_events" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the schedule events.
  If omitted, the provider-level region will be used.

* `event_id` - (Optional, String) Specifies the event ID to filter by.

* `instance_id` - (Optional, String) Specifies the instance ID to filter by.

* `status` - (Optional, String) Specifies the event status to filter by.
  The valid values are as follows:
  + **WAITING**: Waiting.
  + **INQUIRING**: To be authorized.
  + **SCHEDULED**: To be executed.
  + **EXECUTING**: Executing.
  + **COMPLETED**: Completed.
  + **FAILED**: Failed.
  + **CANCELED**: Canceled.

* `type` - (Optional, String) Specifies the event type to filter by.
  The valid value is **RESTAT_NODE**.

* `level` - (Optional, String) Specifies the event level to filter by.
  The valid values are as follows:
  + **CRITICAL**: Critical.
  + **MAJOR**: Major.
  + **MINOR**: Minor.
  + **INFO**: Information.

* `sort_field` - (Optional, String) Specifies the field to sort by.
  The valid values are **planned_execution_time**, **created_time**, **latest_execution_time**.

* `order` - (Optional, String) Specifies the sort order.
  The valid values are **DESC** (descending) and **ASC** (ascending). Defaults to **DESC**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - The total number of events.

* `inquiring_count` - The number of events waiting for authorization.

* `schedule_count` - The number of events to be executed.

* `executing_count` - The number of events being executed.

* `failed_count` - The number of events that failed to execute.

* `events` - The list of schedule events.
  The [events](#events_attr) structure is documented below.

<a name="events_attr"></a>
The `events` block supports:

* `id` - The event ID.

* `instance_id` - The instance ID.

* `instance_name` - The instance name.

* `db_type` - The database type. The valid values are **mysql**, **postgresql**, **sqlserver**.

* `created_time` - The time when the event was created.

* `update_time` - The time when the event was updated.

* `type` - The event type.

* `impact` - The impact of the event on the system.

* `status` - The event status.

* `reason` - The reason for the event.

* `level` - The severity level of the event.

* `execute_time` - The time when the event was executed.

* `latest_execution_time` - The time when the event was last executed.

* `execution_time_window` - The execution time window of the event.
  The [execution_time_window](#events_execution_time_window_attr) structure is documented below.

<a name="events_execution_time_window_attr"></a>
The `execution_time_window` block supports:

* `planned_execution_time` - The planned execution time.

* `start_time` - The start time of the execution window.

* `end_time` - The end time of the execution window.
