---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_ticket_add"
description: |-
  Manages a COC ticket add resource within HuaweiCloud.
---

# huaweicloud_coc_ticket_add

Manages a COC ticket add resource within HuaweiCloud.

~> Deleting ticket add resource is not supported, it will only be removed from the state.

## Example Usage

### add issue ticket

```hcl
variable "title" {}
variable "region" {}
variable "application_id" {}
variable "source_id" {}
variable "found_time" {}
variable "user_id" {}

resource "huaweicloud_coc_ticket_add" "test" {
  ticket_type              = "issues_mgmt"
  title                    = var.title
  description              = "this is description"
  enterprise_project_id    = "0"
  issue_ticket_type        = "issues_type_1000"
  virtual_schedule_type    = "issues_mgmt_virtual_schedule_type_2000"
  regions                  = var.region
  level                    = "issues_level_4000"
  root_cause_cloud_service = var.application_id
  source                   = "issues_mgmt_associated_type_1000"
  source_id                = var.source_id
  found_time               = var.found_time
  issue_contact_person     = var.user_id
}
```

### add change ticket

```hcl
variable "title" {}
variable "change_guide" {}
variable "application_id" {}
variable "application_name" {}
variable "schedule_scene" {}
variable "schedule_role" {}
variable "schedule_role_name" {}
variable "region" {}
variable "expected_start_time" {}
variable "expected_end_time" {}
variable "user_id" {}

resource "huaweicloud_coc_ticket_add" "test" {
  ticket_type              = "change"
  title                    = var.title
  change_notes             = "this is change description"
  change_type              = "change_type_urgentu"
  level                    = "change_level_040"
  change_scheme            = "change_scheme_guides"
  change_guides            = var.change_guide
  change_scene_id          = "GOCMLL01001"
  current_cloud_service_id = var.application_id
  schedule_scenes          = var.schedule_scene
  schedule_roles           = var.schedule_role
  schedule_roles_name      = var.schedule_role_name
  approve_type             = "or_sign"
  location_id              = var.region
  is_start_process         = true
  sub_tickets {
    target_type = "change_scope"
    app_name    = var.application_name
    region      = var.region
  }
  sub_tickets {
    target_type         = "child_ticket"
    target_value        = var.region
    expected_start_time = var.expected_start_time
    expected_end_time   = var.expected_end_time
    executors           = var.user_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `ticket_type` - (Required, String, NonUpdatable) Specifies the type of work order to be created.
  Values can be as follows:
  + **change**: Change order.
  + **issues_mgmt**: Issue order.

* `title` - (Required, String, NonUpdatable) Specifies the ticket title.

* `change_notes` - (Optional, String, NonUpdatable) Specifies the description for the change order.

  -> when `ticket_type` is **change**, this field is required.

* `description` - (Optional, String, NonUpdatable) Specifies the description for the issue ticket.

  -> when `ticket_type` is **issues_mgmt**, this field is required.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.

  -> when `ticket_type` is **issues_mgmt**, this field is required.

* `change_type` - (Optional, String, NonUpdatable) Specifies the change type.
  Values can be as follows:
  + **change_type_conventional**: Conventional change.
  + **change_type_urgentu**: Urgent change.

  -> when `ticket_type` is **change**, this field is required.

* `level` - (Optional, String, NonUpdatable) Specifies the work order level.
  + When the `ticket_type` is **change**, the value can be as follows:
      - **change_level_010**: Level A change.
      - **change_level_020**: Level B change.
      - **change_level_030**: Level C change.
      - **change_level_040**: Level D change.
  + When the `ticket_type` is **issues_mgmt**, the value can be as follows:
      - **issues_level_1000**: Fatal.
      - **issues_level_2000**: Severe.
      - **issues_level_3000**: General.
      - **issues_level_4000**: Info.

* `issue_ticket_type` - (Optional, String, NonUpdatable) Specifies the issue ticket type.

  -> When `ticket_type` is **issues_mgmt**, this field is required.

* `change_scheme` - (Optional, String, NonUpdatable) Specifies the task type.
  Values can be as follows:
  + **change_scheme_runbook**: Document.
  + **change_scheme_guides**: Change guide.

  -> When `ticket_type` is **change**, this field is required.

* `change_guides` - (Optional, String, NonUpdatable) Specifies the change guide ID. Multiple change guide are separated
  by commas.

  -> When `ticket_type` is **change** and `change_scheme` is **change_scheme_guides**, this field is required.

* `commit_upload_attachment` - (Optional, String, NonUpdatable) Specifies the issue ticket attachment ID.

  -> When `ticket_type` is **issues_mgmt**, this field is optional.

* `regions` - (Optional, String, NonUpdatable) Specifies the region to which the issue belongs. Multiple regions are
  separated by commas.

  -> When `ticket_type` is **issues_mgmt**, this field is optional.

* `change_scene_id` - (Optional, String, NonUpdatable) Specifies the change scenario.

  -> When `ticket_type` is **change**, this field is required.

* `current_cloud_service_id` - (Optional, String, NonUpdatable) Specifies the change service ID.

  -> When `ticket_type` is **change**, this field is required.

* `root_cause_cloud_service` - (Optional, String, NonUpdatable) Specifies the root cause service ID of the problem.

  -> When `ticket_type` is **issues_mgmt**, this field is optional.

* `source` - (Optional, String, NonUpdatable) Specifies the source of the ticket.
  Values can be as follows:
  + **issues_mgmt_associated_type_1000**: Incident.
  + **issues_mgmt_associated_type_2000**: Alarm.
  + **issues_mgmt_associated_type_3000**: War room.
  + **issues_mgmt_associated_type_4000**: O&M proactive discovery.

  -> When `ticket_type` is **issues_mgmt**, this field is optional.

* `source_id` - (Optional, String, NonUpdatable) Specifies the work order number from which the problem originates.

  -> When `ticket_type` is **issues_mgmt**, this field is optional. And `source` is **issues_mgmt_associated_type_1000**,
  **issues_mgmt_associated_type_2000** or **issues_mgmt_associated_type_3000**, this field is required.

* `found_time` - (Optional, Int, NonUpdatable) Specifies the discovery time, millisecond timestamp.

  -> When `ticket_type` is **issues_mgmt**, this field is optional.

* `virtual_schedule_type` - (Optional, String, NonUpdatable) Specifies the ticket handler type.
  Values can be as follows:
  + **issues_mgmt_virtual_schedule_type_1000**: Schedule.
  + **issues_mgmt_virtual_schedule_type_2000**: Individual.

  -> When `ticket_type` is **issues_mgmt**, this field is required.

* `issue_contact_person` - (Optional, String, NonUpdatable) Specifies the ID of the person responsible for the problem
  ticket.

  -> When `ticket_type` is **issues_mgmt**, this field is optional. And `virtual_schedule_type` is
  **issues_mgmt_virtual_schedule_type_2000**, this field is required.

* `schedule_scenes` - (Optional, String, NonUpdatable) Specifies the scheduling scenario ID. Multiple scenarios are
  separated by commas.

* `schedule_roles` - (Optional, String, NonUpdatable) Specifies the scheduling role ID. Multiple roles are separated by
  commas.

* `schedule_roles_name` - (Optional, String, NonUpdatable) Specifies the scheduling role name. Multiple role names are
  separated by commas.

* `approve_type` - (Optional, String, NonUpdatable) Specifies the approval type. Multiple approve type are separated by
  commas.

  -> When `ticket_type` is **change**, the field `schedule_scenes`, `schedule_roles`, `schedule_roles_name` and
  `approve_type` is required.

* `virtual_schedule_role` - (Optional, String, NonUpdatable) Specifies the scheduling role ID.

  -> When `ticket_type` is **issues_mgmt**, this field is optional.

* `location_id` - (Optional, String, NonUpdatable) Specifies the ID of the region to be changed. Multiple location are
  separated by commas.

  -> When `ticket_type` is **change**, this field is required.

* `plan_task_sub_type` - (Optional, String, NonUpdatable) Specifies the plan subtype.
  Values can be as follows:
  + **CUSTOMIZATION**: Customized job.
  + **COMMUNAL**: Public job.

* `plan_task_id` - (Optional, String, NonUpdatable) Specifies the task ID to be executed.

* `plan_task_name` - (Optional, String, NonUpdatable) Specifies the task name to be executed.

* `plan_task_params` - (Optional, String, NonUpdatable) Specifies the parameter information required to execute the task.

-> When `ticket_type` is **change** and `change_scheme` is **change_scheme_runbook**, the field `plan_task_sub_type`,
`plan_task_id`, `plan_task_id` and `plan_task_id` is required.

* `is_start_process` - (Optional, Bool, NonUpdatable) Specifies whether to start the process.
  + When the value is **false**, the created ticket will be in draft status.
  + When the value is **true**, the created ticket will be in unaccepted status.

* `sub_tickets` - (Optional, List, NonUpdatable) Specifies the information of the change sub-order.

  The [sub_tickets](#sub_tickets_struct) structure is documented below.

  -> When `ticket_type` is **change**, this field is required.

<a name="sub_tickets_struct"></a>
The `sub_tickets` block supports:

* `target_type` - (Optional, String, NonUpdatable) Specifies the target option information.
  Values can be as follows:
  + **change_scope**: Change application.
  + **child_ticket**: Change plan.

* `app_name` - (Optional, String, NonUpdatable) Specifies the change service.

  -> When `target_type` is set to **change_scope**, this field must contain the corresponding change service name in Chinese.

* `region` - (Optional, String, NonUpdatable) Specifies the change region.

  -> When `target_type` is set to **change_scope**, this field must contain the corresponding change region ID.

* `target_value` - (Optional, String, NonUpdatable) Specifies the region ID of the change order to be transferred.

  -> This value is valid only when `target_type` is **child_ticket**.

* `expected_start_time` - (Optional, Int, NonUpdatable) Specifies the timestamp of the change order's scheduled start time.

  -> This value is valid only when `target_type` is **child_ticket**.

* `expected_end_time` - (Optional, Int, NonUpdatable) Specifies the timestamp of the change order's scheduled end time.

  -> This value is valid only when `target_type` is **child_ticket**.

* `executors` - (Optional, String, NonUpdatable) Specifies the person who will implement the change order.

  -> This value is valid only when `target_type` is **child_ticket**.

* `cooperators` - (Optional, String, NonUpdatable) Specifies the person to collaborate on the change order.

  -> This value is valid only when `target_type` is **child_ticket**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `ticket_id` - Indicates the number of the newly created work order.

* `is_start_process_async` - Indicates whether to start the process asynchronously.

* `is_update_null` - Indicates whether to resubmit empty fields.

* `is_return_full_info` - Indicates whether to return all fields.
