---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_issue"
description: |-
  Manages a COC issue resource within HuaweiCloud.
---

# huaweicloud_coc_issue

Manages a COC issue resource within HuaweiCloud.

~> Deleting issue resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "title" {}
variable "ticket_type" {}
variable "virtual_schedule_type" {}
variable "level" {}

resource "huaweicloud_coc_issue" "test" {
  title                 = var.title
  description           = "this is description"
  enterprise_project_id = "0"
  ticket_type           = var.ticket_type
  virtual_schedule_type = var.virtual_schedule_type
  level                 = var.level
}
```

## Argument Reference

The following arguments are supported:

* `title` - (Required, String, NonUpdatable) Specifies the ticket title.

* `description` - (Required, String, NonUpdatable) Specifies the description for the ticket.

* `enterprise_project_id` - (Required, String, NonUpdatable) Specifies the enterprise project ID.

* `ticket_type` - (Required, String, NonUpdatable) Specifies the ticket type.

* `virtual_schedule_type` - (Required, String, NonUpdatable) Specifies the ticket handler type.
  The valid values are as follows:
  + **issues_mgmt_virtual_schedule_type_1000**: Schedule.
  + **issues_mgmt_virtual_schedule_type_2000**: Individual.

* `commit_upload_attachment` - (Optional, String, NonUpdatable) Specifies the issue ticket attachment ID.

* `regions` - (Optional, List, NonUpdatable) Specifies the region to which the issue belongs.

* `level` - (Optional, String, NonUpdatable) Specifies the work order level.
  The valid values are as follows:
  + **issues_level_1000**: Fatal.
  + **issues_level_2000**: Severe.
  + **issues_level_3000**: General.
  + **issues_level_4000**: Info.

* `root_cause_cloud_service` - (Optional, String, NonUpdatable) Specifies the root cause service ID of the problem.

* `source` - (Optional, String, NonUpdatable) Specifies the source of the ticket.
  The valid values are as follows:
  + **issues_mgmt_associated_type_1000**: Incident.
  + **issues_mgmt_associated_type_2000**: Alarm.
  + **issues_mgmt_associated_type_3000**: War room.
  + **issues_mgmt_associated_type_4000**: O&M proactive discovery.

* `source_id` - (Optional, String, NonUpdatable) Specifies the work order number from which the problem originates.

* `found_time` - (Optional, Int, NonUpdatable) Specifies the discovery time, millisecond timestamp.

* `issue_contact_person` - (Optional, String, NonUpdatable) Specifies the ID of the person responsible for the issue.

* `schedule_scenes` - (Optional, String, NonUpdatable) Specifies the scheduling scenario ID.

* `virtual_schedule_role` - (Optional, String, NonUpdatable) Specifies the scheduling role ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `is_start_process_async` - Indicates whether to start the process asynchronously.

* `is_update_null` - Indicates whether to resubmit empty fields.

* `is_return_full_info` - Indicates whether to return all fields.

* `is_start_process` - Indicates whether to start the process.

* `ticket_id` - Indicates the number of the newly created work order.

* `issue_correlation_sla` - Indicates the SLA associated with the problem.

* `root_cause_type` - Indicates the classification of single root causes of problems.

* `current_cloud_service_id` - Indicates the problem ticket service.

* `issue_version` - Indicates the version number where the problem was found.

* `root_cause_comment` - Indicates the single root cause analysis of a problem.

* `solution` - Indicates the problem ticket resolution.

* `regions_search` - Indicates the value representing a ticket region search.

* `level_approve_config` - Indicates the ticket level approval configuration.

* `suspension_approve_config` - Indicates the configuration of pending approval for a problem ticket.

* `handle_time` - Indicates the time it takes to handle a problem ticket.

* `is_common_issue` - Indicates whether it is a common problem.

* `is_need_change` - Indicates whether the issue ticket requires changes.

* `is_enable_suspension` - Indicates whether the ticket is open or pending.

* `creator` - Indicates the creator of the problem ticket.

* `operator` - Indicates the operator of the problem ticket.

* `real_ticket_id` - Indicates the problem ticket number.

* `assignee` - Indicates the current person responsible for the issue ticket.

* `participator` - Indicates the ticket participant.

* `work_flow_status` - Indicates the status of the problem ticket.

* `engine_error_msg` - Indicates the process status.

* `baseline_status` - Indicates the baseline state.

* `phase` - Indicates the current stage information of the issue ticket.

* `sub_tickets` - Indicates the change sub-order information.

  The [sub_tickets](#sub_tickets_struct) structure is documented below.

* `enum_data_list` - Indicates the enumerated list representing issue associations.

  The [enum_data_list](#enum_data_list_struct) structure is documented below.

* `meta_data_version` - Indicates the application version that caused the problem.

* `update_time` - Indicates the update time.

* `create_time` - Indicates the creation time.

* `is_deleted` - Indicates whether the work order is deleted.

* `ticket_type_id` - Indicates the work order type.

* `form_info` - Indicates the action information.

<a name="sub_tickets_struct"></a>
The `sub_tickets` block supports:

* `is_deleted` - Indicates whether it has been deleted.

* `id` - Indicates the work order ID.

* `main_ticket_id` - Indicates the work order primary key ID.

* `parent_ticket_id` - Indicates the parent work order ID.

* `ticket_id` - Indicates the problem ticket ID.

* `real_ticket_id` - Indicates the problem ticket number.

* `ticket_path` - Indicates the work order path.

* `target_value` - Indicates the region information.

* `target_type` - Indicates the sub-order type.

* `create_time` - Indicates the time when the work order was created.

* `update_time` - Indicates the time when the work order was updated.

* `creator` - Indicates the ID of the work order creator.

* `operator` - Indicates the ID of the work order operator.

<a name="enum_data_list_struct"></a>
The `enum_data_list` block supports:

* `is_deleted` - Indicates whether it has been deleted.

* `match_type` - Indicates the matching enumeration type.

* `ticket_id` - Indicates the current work order ID.

* `real_ticket_id` - Indicates the work order number.

* `name_zh` - Indicates the Chinese name.

* `name_en` - Indicates the English name.

* `biz_id` - Indicates the unique id corresponding to the enumeration value.

* `prop_id` - Indicates the type corresponding to the current enumeration value.

* `model_id` - Indicates the model ID corresponding to different background applications.

## Import

The COC issue can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_coc_issue.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `handle_time`.
It is generally recommended running `terraform plan` after importing a issue.
You can then decide if changes should be applied to the issue, or the resource definition should be updated to
align with the issue. Also you can ignore changes as below.

```hcl
resource "huaweicloud_coc_issue" "test" {
    ...

  lifecycle {
    ignore_changes = [
      handle_time,
    ]
  }
}
```
