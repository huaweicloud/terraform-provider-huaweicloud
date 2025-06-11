---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_war_room"
description: |-
  Manages a COC war room resource within HuaweiCloud.
---

# huaweicloud_coc_war_room

Manages a COC war room resource within HuaweiCloud.

~> Deleting war room resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
var "war_room_name" {}
var "application_id" {}
var "incident_number" {}
var "role_id" {}
var "scene_id" {}
var "user_id" {}

resource "huaweicloud_coc_war_room" "test" {
  war_room_name         = var.war_room_name
  application_id_list   = [var.application_id]
  incident_number       = var.incident_number
  war_room_admin        = var.user_id
  enterprise_project_id = "0"

  schedule_group {
    role_id  = var.role_id
    scene_id = var.scene_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `war_room_name` - (Required, String, NonUpdatable) Specifies the war room name.

* `application_id_list` - (Required, List, NonUpdatable) Specifies the ID list of the affected applications.

* `incident_number` - (Required, String, NonUpdatable) Specifies the incident ticket ID.

* `schedule_group` - (Required, List, NonUpdatable) Specifies the scheduling group information.
  The [schedule_group](#block--schedule_group) structure is documented below.

* `war_room_admin` - (Required, String, NonUpdatable) Specifies the war room administrator user ID.

* `enterprise_project_id` - (Required, String, NonUpdatable) Specifies the enterprise project ID.

* `description` - (Optional, String, NonUpdatable) Specifies the war room description.

* `region_code_list` - (Optional, List, NonUpdatable) Specifies the ID list of regions.

* `participant` - (Optional, List, NonUpdatable) Specifies the user ID list of participants.

* `application_names` - (Optional, List, NonUpdatable) Specifies the application name list.

* `region_names` - (Optional, List, NonUpdatable) Specifies the region names.

* `notification_type` - (Optional, String, NonUpdatable) Specifies the group creation mode.
  Values can be as follows:
  + **WECHAT**: WeChat.
  + **DING_TALK**: DingTalk.
  + **LARK**: Lark.
  + **NULL_GROUP**: No notification is sent to the group.

<a name="block--schedule_group"></a>
The `schedule_group` block supports:

* `role_id` - (Required, String, NonUpdatable) Specifies the role ID.

* `scene_id` - (Required, String, NonUpdatable) Specifies the scenario ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `war_room_id` - Indicates the war room ID.

* `recover_member` - Indicates the members of recovery.

* `recover_leader` - Indicates the primary recovery owner.

* `incident` - Indicates the incident.

  The [incident](#incident_struct) structure is documented below.

* `source` - Indicates the incident source.

* `change_num` - Indicates the change ticket number.

* `occur_time` - Indicates the occurrence start time.

* `recover_time` - Indicates the fault recovery time.

* `fault_cause` - Indicates the fault cause.

* `create_time` - Indicates the creation time.

* `first_report_time` - Indicates the first notification time.

* `recovery_notification_time` - Indicates the recovery notification time.

* `fault_impact` - Indicates the impact of the fault.

* `circular_level` - Indicates the notification level. The notification level is the same as the incident level in the
  tenant zone.

* `war_room_status` - Indicates the war room status.

  The [war_room_status](#war_room_status_struct) structure is documented below.

* `processing_duration` - Indicates the handling duration, the unit is minutes.

* `restoration_duration` - Indicates the recovery duration, the unit is minutes.

<a name="incident_struct"></a>
The `incident` block supports:

* `id` - Indicates the incident primary key.

* `incident_id` - Indicates the incident ID.

* `is_change_event` - Indicates whether the incident is a change incident.

* `failure_level` - Indicates the incident level.

* `incident_url` - Indicates the incident URL.

<a name="war_room_status_struct"></a>
The `war_room_status` block supports:

* `id` - Indicates the war room status enumeration value ID .

* `name_zh` - Indicates the Chinese name of the war room status enumeration value.

* `name_en` - Indicates the English name of the war room status enumeration value.

* `type` - Indicates the war room status enumeration type.

## Import

The COC war room can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_coc_war_room.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `schedule_group`, `participant` and `notification_type`.
It is generally recommended running `terraform plan`  after importing a war room.
You can then decide if changes should be applied to the war room, or the resource definition should be updated to
align with the war room. Also you can ignore changes as below.

```hcl
resource "huaweicloud_coc_war_room" "test" {
    ...

  lifecycle {
    ignore_changes = [
      schedule_group, participant, notification_type,
    ]
  }
}
```
