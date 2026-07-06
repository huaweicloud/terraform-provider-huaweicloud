---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_events"
description: |-
  Use this data source to get the list of DSC security events within HuaweiCloud.
---

# huaweicloud_dsc_events

Use this data source to get the list of DSC security events within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dsc_events" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `event_level` - (Optional, String) Specifies the event risk level.

* `event_name` - (Optional, String) Specifies the event name.

* `event_status` - (Optional, String) Specifies the event handling status.

* `end_time` - (Optional, Int) Specifies the end time of the query (timestamp).

* `start_time` - (Optional, Int) Specifies the start time of the query (timestamp).

* `responsible_person` - (Optional, String) Specifies the responsible person for the event.

* `source_name` - (Optional, String) Specifies the event source name.

* `source_type` - (Optional, String) Specifies the event source type.

* `verification_status` - (Optional, String) Specifies the verification status.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `events` - The event information list.

  The [events](#events_struct) structure is documented below.

<a name="events_struct"></a>
The `events` block supports:

* `id` - The event ID.

* `event_name` - The event name.

* `event_level` - The event risk level.

* `event_status` - The event handling status.

* `event_type` - The event type.

* `close_reason` - The reason for closing the event.

* `create_time` - The event creation time.

* `description` - The event description.

* `disposal_suggestion` - The disposal suggestion.

* `domain_id` - The domain ID.

* `occur_time` - The event occurrence time.

* `project_id` - The project ID.

* `scheduled_close_time` - The scheduled close time.

* `source_module` - The event source module.

* `verification_status` - The verification status.

* `affected_asset` - The affected assets.

* `responsible_person` - The responsible person information.

  The [responsible_person](#events_responsible_person_struct) structure is documented below.

* `source_instance` - The event source instance information.

  The [source_instance](#events_source_instance_struct) structure is documented below.

* `related_alarm_list` - The related alarm information list.

  The [related_alarm_list](#events_related_alarm_list_struct) structure is documented below.

<a name="events_responsible_person_struct"></a>
The `responsible_person` block supports:

* `user_id` - The user ID.

* `user_name` - The user name.

<a name="events_source_instance_struct"></a>
The `source_instance` block supports:

* `instance_id` - The instance ID.

* `instance_name` - The instance name.

<a name="events_related_alarm_list_struct"></a>
The `related_alarm_list` block supports:

* `id` - The alarm ID.

* `alarm_name` - The alarm name.

* `alarm_level` - The alarm level.

* `alarm_status` - The alarm status.

* `alarm_type` - The alarm type.

* `close_reason` - The reason for closing the alarm.

* `create_time` - The alarm creation time.

* `description` - The alarm description.

* `disposal_suggestion` - The disposal suggestion.

* `domain_id` - The domain ID.

* `occur_time` - The alarm occurrence time.

* `project_id` - The project ID.

* `source_module` - The alarm source module.

* `source_sub_type` - The alarm source sub type.

* `source_type` - The alarm source type.

* `verification_status` - The verification status.

* `affected_asset` - The affected assets.

* `responsible_person` - The responsible person information.

  The [responsible_person](#events_related_alarm_responsible_person_struct) structure is documented below.

* `source_instance` - The alarm source instance information.

  The [source_instance](#events_related_alarm_source_instance_struct) structure is documented below.

<a name="events_related_alarm_responsible_person_struct"></a>
The `responsible_person` block supports:

* `user_id` - The user ID.

* `user_name` - The user name.

<a name="events_related_alarm_source_instance_struct"></a>
The `source_instance` block supports:

* `instance_id` - The instance ID.

* `instance_name` - The instance name.
