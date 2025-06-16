---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_war_rooms"
description: |-
  Use this data source to get the list of COC war rooms.
---

# huaweicloud_coc_war_rooms

Use this data source to get the list of COC war rooms.

## Example Usage

```hcl
data "huaweicloud_coc_war_rooms" "test" {}
```

## Argument Reference

The following arguments are supported:

* `incident_num` - (Optional, String) Specifies the incident ticket number.

* `title` - (Optional, String) Specifies the war room name.

* `region_code_list` - (Optional, List) Specifies the regions.

* `incident_levels` - (Optional, List) Specifies the incident level.
  Values can be as follows:
  + **level_10**: P1.
  + **level_20**: P2.
  + **level_30**: P3.
  + **level_40**: P4.
  + **level_50**: P5.

* `impacted_application_ids` - (Optional, List) Specifies the ID of the affected application.

* `admin` - (Optional, List) Specifies the war room administrator user ID.

* `status` - (Optional, List) Specifies the war room status.
  Values can be as follows:
  + **1**: Start war room.
  + **3**: Fault definition.
  + **7**: The fault has been recovered.
  + **20**: Close war room.

* `triggered_start_time` - (Optional, Int) Specifies the trigger start time of the war room. The default value is 30
  days before the start time.

* `triggered_end_time` - (Optional, Int) Specifies the trigger end time of the war room. The default value is current
  time.

* `occur_start_time` - (Optional, Int) Specifies the occurrence start time.

* `occur_end_time` - (Optional, Int) Specifies the occurrence end time.

* `recover_start_time` - (Optional, Int) Specifies the recovery start time.

* `recover_end_time` - (Optional, Int) Specifies the recovery end time.

* `notification_level` - (Optional, List) Specifies the notification level.
  Values can be as follows:
  + **level_10**: P1.
  + **level_20**: P2.
  + **level_30**: P3.
  + **level_40**: P4.
  + **level_50**: P5.

* `enterprise_project_ids` - (Optional, List) Specifies the enterprise project ID.

* `war_room_num` - (Optional, String) Specifies the war room ticket number.

* `statistic_flag` - (Optional, Bool) Specifies whether to collect statistics. If the value is **false**, basic
  information is returned. If the value is **true**, only the statistics result, including `total_num`, `running_num`,
  and `closed_num`, is returned. The default value is **false**.

* `current_users` - (Optional, List) Specifies the current user ID.

* `war_room_nums` - (Optional, List) Specifies the war room ticket number. When this filter is present, other filter
  conditions are ignored.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `list` - Indicates the war room information.

  The [list](#data_list_struct) structure is documented below.

* `running_num` - Indicates the total number of war rooms in progress.

* `closed_num` - Indicates the total number of closed war rooms.

* `total_num` - Indicates the total number of war rooms.

<a name="data_list_struct"></a>
The `list` block supports:

* `id` - Indicates the primary key.

* `title` - Indicates the title.

* `admin` - Indicates the war room administrator user ID.

* `recover_member` - Indicates the members of recovery.

* `recover_leader` - Indicates the primary recovery owner.

* `incident` - Indicates the incident.

  The [incident](#list_incident_struct) structure is documented below.

* `source` - Indicates the incident source.

* `regions` - Indicates the affected regions.

  The [regions](#list_regions_struct) structure is documented below.

* `change_num` - Indicates the change ticket number.

* `occur_time` - Indicates the occurrence start time.

* `recover_time` - Indicates the fault recovery time.

* `fault_cause` - Indicates the fault cause.

* `create_time` - Indicates the creation time.

* `first_report_time` - Indicates the first notification time.

* `recovery_notification_time` - Indicates the recovery notification time.

* `fault_impact` - Indicates the impact of the fault.

* `description` - Indicates the war room description.

* `circular_level` - Indicates the notification level. The notification level is the same as the incident level in the
  tenant zone.

* `war_room_status` - Indicates the war room status.

  The [war_room_status](#list_war_room_status_struct) structure is documented below.

* `impacted_application` - Indicates the affected applications.

  The [impacted_application](#list_impacted_application_struct) structure is documented below.

* `processing_duration` - Indicates the handling duration, the unit is minutes.

* `restoration_duration` - Indicates the recovery duration, the unit is minutes.

* `war_room_num` - Indicates the war room ticket number.

* `enterprise_project_id` - Indicates the enterprise project ID.

<a name="list_incident_struct"></a>
The `incident` block supports:

* `id` - Indicates the incident primary key.

* `incident_id` - Indicates the incident ID.

* `is_change_event` - Indicates whether the incident is a change incident.

* `failure_level` - Indicates the incident level.

* `incident_url` - Indicates the incident URL.

<a name="list_regions_struct"></a>
The `regions` block supports:

* `code` - Indicates the region primary key.

* `name` - Indicates the region name.

<a name="list_war_room_status_struct"></a>
The `war_room_status` block supports:

* `id` - Indicates the war room status enumeration value ID .

* `name_zh` - Indicates the Chinese name of the war room status enumeration value.

* `name_en` - Indicates the English name of the war room status enumeration value.

* `type` - Indicates the war room status enumeration type.

<a name="list_impacted_application_struct"></a>
The `impacted_application` block supports:

* `id` - Indicates the affected application primary key.

* `name` - Indicates the affected application name.
