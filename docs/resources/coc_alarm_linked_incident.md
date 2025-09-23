---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_alarm_linked_incident"
description: |-
  Manages a COC alarm linked incident resource within HuaweiCloud.
---

# huaweicloud_coc_alarm_linked_incident

Manages a COC alarm linked incident resource within HuaweiCloud.

~> Deleting alarm linked incident resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "alarm_id" {}
variable "application_id" {}
variable "title" {}

resource "huaweicloud_coc_alarm_linked_incident" "test" {
  alarm_ids                = var.alarm_id
  current_cloud_service_id = var.application_id
  description              = "alarm to incident"
  is_service_interrupt     = false
  level_id                 = "level_50"
  mtm_type                 = "inc_type_p_change_issues"
  title                    = var.title
}
```

## Argument Reference

The following arguments are supported:

* `alarm_ids` - (Required, String, NonUpdatable) Specifies the list of alarm IDs, separated by commas.

* `current_cloud_service_id` - (Required, String, NonUpdatable) Specifies the fault application ID.

* `description` - (Required, String, NonUpdatable) Specifies the event description.

* `is_service_interrupt` - (Required, Bool, NonUpdatable) Specifies the service is interrupted. The default value is **false**.
  Valid values are as follows:
  + **false**: The service is not interrupted.
  + **true**: The service is interrupted.

* `level_id` - (Required, String, NonUpdatable) Specifies the event severity. The default value is **level_50**.
  Valid values are as follows:
  + **level_10**: P1 event.
  + **level_20**: P2 event.
  + **level_30**: P3 event.
  + **level_40**: P4 event.
  + **level_50**: P5 event.

* `mtm_type` - (Required, String, NonUpdatable) Specifies the event category.

* `title` - (Required, String, NonUpdatable) Specifies the event name.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.

* `assignee` - (Optional, String, NonUpdatable) Specifies the assignee.

* `assignee_role` - (Optional, String, NonUpdatable) Specifies the scheduling role.

* `assignee_scene` - (Optional, String, NonUpdatable) Specifies the scheduling scenario.

* `attachment` - (Optional, String, NonUpdatable) Specifies the list of attachments.

* `is_change_event` - (Optional, Bool, NonUpdatable) Specifies whether the event is changed. The default value is **false**.
  Valid values are as follows:
  + **false**: The event is not interrupted.
  + **true**: The event is interrupted.

* `mtm_region` - (Optional, String, NonUpdatable) Specifies the region ID.

* `source_id` - (Optional, String, NonUpdatable) Specifies how the event was created.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
