---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_ticket_action"
description: |-
  Manages a COC ticket action resource within HuaweiCloud.
---

# huaweicloud_coc_ticket_action

Manages a COC ticket action resource within HuaweiCloud.

~> Deleting ticket action resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "ticket_id" {}

resource "huaweicloud_coc_ticket_action" "test" {
  ticket_type = "issues_mgmt"
  ticket_id   = var.ticket_id
  action      = "gocm_issues_accepte"
}
```

## Argument Reference

The following arguments are supported:

* `ticket_type` - (Required, String, NonUpdatable) Specifies the type of work order that needs to be operated.
  The valid values are as follows:
  + **change**: Change order.
  + **issues_mgmt**: Problem ticket.

* `ticket_id` - (Required, String, NonUpdatable) Specifies the work order number that requires operation.

* `user_id` - (Optional, String, NonUpdatable) Specifies the user ID performing the current operation.
  The default value obtained from the authentication information.

* `task_id` - (Optional, String, NonUpdatable) Specifies the task ID corresponding to the current operation.

* `action` - (Optional, String, NonUpdatable) Specifies the action information that needs to be performed.
  The valid values are as follows:
  + **gocm_issues_accepte**: Issue acceptance.
  + **ugrading_and_downgrading**: Issue upgrade/downgrading.
  + **being_handled_initiate_suspend**: Issue pending.
  + **agreed**: Issue pending approval or change order approval.
  + **rejected**: Change order approval rejected.
  + **suspend_resume**: Issue pending resumed.
  + **gocm_issues_handling_forwarding**: Issue forwarding responsible.
  + **gocm_issues_unaccepted_reject**: Issue rejection.
  + **gocm_issues_undo_create**: Issue cancellation.
  + **submitting_question**: Issue solution submission.
  + **cancel_process_action**: Change order cancellation.

* `params` - (Optional, String, NonUpdatable) Specifies the specific action information for executing the operation.
  The following action can be applied, and different action contain different fields:
  + gocm_issues_accepte: Issue acceptance. No fields.
  + ugrading_and_downgrading: Issue upgrade/downgrading.
      - **up_down_updated_level**: Upgrade or downgrade a work order level.
      - **up_down_grade_reason_comment**: Upgrade or downgrade a work order comment.
      - **level**: Work order level.
  + being_handled_initiate_suspend: Issue pending.
      - **estimated_recovery_time**: Estimated recovery time, millisecond timestamp.
      - **suspend_description**: Suspend description.
  + agreed: Issue pending approval or change order approval.
      - **conclusion**: Conclusion. The value can be `true` or `false`.
      - **approve_note**: Approved description.
  + rejected: Change order approval rejected.
      - **comment**: Comment.
  + suspend_resume: Issue pending resumed. No fields.
  + gocm_issues_handling_forwarding: Issue forwarding responsible.
      - **explain**: Reason.
      - **forwarding_assigne**: Forwarding assignee.
      - **virtual_schedule_type**: Schedule type.
  + gocm_issues_unaccepted_reject: Issue rejection.
      - **unaccepted_rejection_description**: Unaccepted rejection description.
  + gocm_issues_undo_create: Issue cancellation.
      - **revocation_reason**: Reason.
  + submitting_question: Issue solution submission.
      - **root_cause_cloud_service**: The root cause cloud service.
      - **is_common_issue**: Whether it contains common problems. The value can be `true` or `false`.
      - **issue_version**: The version number in which the problem was found.
      - **regions**: Region information.
      - **root_cause_type**: Root cause type.
      - **root_cause_comment**: Root cause comment.
      - **solution**: Solution.
      - **is_need_change**: Whether changes are required. The value can be `true` or `false`.
      - **sub_tickets**: The list of change sub-order information.
          + **affected_region**: Affected region.
          + **target_type**: Sub-order type.
          + **target_value**: Region information.
  + cancel_process_action: Change order cancellation.
      - **comment**: Comment.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `ticket_id`.
