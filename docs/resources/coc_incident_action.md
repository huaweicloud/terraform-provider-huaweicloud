---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_incident_action"
description: |-
  Manages a COC incident action resource within HuaweiCloud.
---

# huaweicloud_coc_incident_action

Manages a COC incident action resource within HuaweiCloud.

~> Deleting incident action resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "incident_id" {}
variable "task_id" {}

resource "huaweicloud_coc_incident_action" "test" {
  incident_id = var.incident_id
  task_id     = var.task_id
  action      = "accepted"
}
```

## Argument Reference

The following arguments are supported:

* `incident_id` - (Required, String, NonUpdatable) Specifies the event ticket number.

* `task_id` - (Required, String, NonUpdatable) Specifies the task ID.
  The value depend on the `task_id` returned by the `huaweicloud_coc_incident_tasks`.

* `action` - (Required, String, NonUpdatable) Specifies the identification of different operation types.
  The valid values are as follows, from the `key` returned by the `huaweicloud_coc_incident_tasks`:
  + **unAcceptedForward**: Forwarded to the responsible person in the unaccepted state.
  + **acceptedForward**: Forward to the responsible person after acceptance.
  + **accepted**: Accepted.
  + **addRemark**: Add remark.
  + **commitSolution**: Solution.
  + **confirm**: Verification closed.
  + **incidentPause**: Request for incident pause.
  + **agreed**: Agree under suspended approval status.
  + **rejected**: Rejected under suspended approval status.
  + **recovery**: Recovered after suspension.
  + **changeLevel**: Upgrade/downgrade request.
  + **agreed**: Agree under the upgrade/downgrade approval status.
  + **rejected**: Rejected under the upgrade/downgrade approval status.
  + **rejected**: Rejected in unaccepted state.
  + **Reopen**: Reopen in rejected state.
  + **agreed**: Closing the order in the rejected state.

* `params` - (Optional, Map, NonUpdatable) Specifies the parameters.
  The following `action` can be applied, and different actions correspond to different parameters and their required items:
  + unAcceptedForward: Forwarded to the responsible person in the unaccepted state.
      - **virtual_schedule_type**: Scheduling scenario.
      - **virtual_schedule_role**: Scheduling roles.
      - **virtual_send_assignee**: Responsible person for forwarding.
      - **virtual_send_comment**: Remarks information.
      - **virtual_current_location_info**: Positioning status at the current stage.
  + acceptedForward: Forward to the responsible person after acceptance.
      - **virtual_schedule_type**: Scheduling scenario.
      - **virtual_schedule_role**: Scheduling roles.
      - **virtual_send_assignee**: Responsible person for forwarding.
      - **virtual_send_comment**: Remarks information.
      - **virtual_current_location_info**: Positioning status at the current stage.
  + accepted: Accepted. No value needs to be passed.
  + addRemark: Add remark.
      - **note**: Remarks information.
  + commitSolution: Solution.
      - **mtm_type**: Question type.
      - **is_service_interrupt**: Whether the business is interrupted.
      - **start_time**: Failure time. The value is required when `is_service_interrupt` is **true**.
      - **fault_recovery_time**: Fault recovery time. The value is required when `is_service_interrupt` is **true**.
      - **cause**: Cause of failure.
      - **solution**: Solution.
      - **resolve_attachments**: Attachment ID.
  + confirm: Verification closed.
      - **virtual_confirm_result**: Solved or not.
      - **virtual_confirm_comment**: Remarks. The value is required for event approval.
  + incidentPause: Request for incident pause.
      - **pause_end_time**: Suspension deadline.
      - **pause_reason**: Reason for suspension.
  + agreed: Agree under suspended approval status.
      - **pause_approve_conclusion**: Whether it has passed the review. The value is **true**.
      - **note**: Remarks information.
  + rejected: Rejected under suspended approval status.
      - **pause_approve_conclusion**: Whether it has passed the review. The value is **false**.
      - **note**: Remarks information.
  + recovery: Recovered after suspension. No value needs to be passed.
  + changeLevel: Upgrade/downgrade request.
      - **virtual_target_level**: Target event level.
      - **virtual_change_level_comment**: Notes added during the upgrade or downgrade of an event ticket.
  + agreed: Agree under the upgrade/downgrade approval status.
      - **conclusion**: Whether to agree or not. The value is **true**.
      - **note**: Remarks information.
  + rejected: Rejected under the upgrade/downgrade approval status.
  + agreed: Agree under the upgrade/downgrade approval status.
      - **conclusion**: Whether to agree or not. The value is **false**.
      - **note**: Remarks information.
  + rejected: Rejected in unaccepted state.
      - **virtual_confirm_comment**: Remarks information.
  + ReOpen: Reopen in rejected state.
      - **mtm_region**: Region.
      - **enterprise_project_id**: The enterprise project ID bound to the event ticket.
      - **current_cloud_service_id**: Cloud service ID. It can be obtained through the query applications.
      - **level_id**: Event level.
        For details, see [level_id](https://support.huaweicloud.com/api-coc/coc_api_04_03_001_006.html#coc_api_04_03_001_006__section289718103710).
      - **is_service_interrupt**: Whether the service is interrupted.
      - **mtm_type**: Question category.
      - **title**: Title.
      - **description**: Description.
      - **attachments**: Attachment ID.
      - **source_id**: Incident source. If the incident is manually created, set it to **incident_source_manual**.
      - **incident_ownership**: Event attribution.
      - **start_time**: Time when the fault occurred.
      - **assignee**: Incident handler.
      - **assignee_scene**: Scheduling scenario.
      - **assignee_role**: Scheduling roles.
  + agreed: Closing the order in the rejected state.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `incident_id`.
