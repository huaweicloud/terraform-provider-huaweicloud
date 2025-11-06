# huaweicloud_aom_alarm_notified_history

Provides a data source to query alarm notification histories in AOM (Application Operations Management).

## Example Usage

```hcl
data "huaweicloud_aom_alarm_notified_history" "test" {
  time_range      = "-1.-1.60"  # Last 60 minutes
  event_type      = "alarm"
  event_severity  = "Critical"
  limit           = 100
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region where the alarm histories are located.

* `time_range` - (Required, String) The time range for querying alarm histories in the format: startTimeInMillis.endTimeInMillis.durationInMinutes.

* `type` - (Optional, String) The type of alarm to query. Valid values: active_alert, history_alert.

* `event_type` - (Optional, String) The event type to filter alarm histories.

* `event_severity` - (Optional, String) The severity level to filter alarm histories. Valid values: Critical, Major, Minor, Info.

* `limit` - (Optional, Int) The maximum number of alarm histories to return. Defaults to 1000.

* `marker` - (Optional, String) The pagination marker for the next page of results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `events` - The list of alarm notification histories. Structure is documented below.

* `next_marker` - The marker for the next page of results.

The `events` block contains:

* `event_id` - The unique ID of the alarm event.

* `event_name` - The name of the alarm event.

* `event_type` - The type of the alarm event.

* `event_severity` - The severity level of the alarm event.

* `resource_id` - The ID of the resource associated with the alarm.

* `resource_name` - The name of the resource associated with the alarm.

* `start_time` - The start time of the alarm event.

* `end_time` - The end time of the alarm event.

* `detail` - The detailed information of the alarm event.