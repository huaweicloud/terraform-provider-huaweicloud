---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_tasks"
description: |-
  Use this data source to get the list of SecMaster tasks within HuaweiCloud.
---

# huaweicloud_secmaster_tasks

Use this data source to get the list of SecMaster tasks within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_tasks" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `sort_key` - (Optional, String) Specifies the sort key. Supported values are **create_time** and **update_time**.

* `sort_dir` - (Optional, String) Specifies the sort direction. Valid values are:
  + **ASC**: Sort in ascending order.
  + **DESC**: Sort in descending order.

* `note` - (Optional, String) Specifies the note of the task.

* `name` - (Optional, String) Specifies the name of the task.

* `business_type` - (Optional, String) Specifies the business type. Supported values are:
  + **WORKFLOWPUBLISH**: Workflow publish.
  + **WORKFLOWNODEAPPROVE**: Workflow node approve.
  + **PLAYBOOKPUBLISH**: Playbook publish.
  + **PLAYBOOKNODEAPPROVE**: Playbook node approve.

* `creator_name` - (Optional, String) Specifies the creator name of the task.

* `query_type` - (Optional, String) Specifies the query type. Supported values are:
  + **my_todo**: My to-do tasks.
  + **all_handled**: All handled tasks.

* `from_date` - (Optional, String) Specifies the start time in format "yyyy-MM-dd'T'HH:mm:ss.SSS'Z'Z".

* `to_date` - (Optional, String) Specifies the end time in format "yyyy-MM-dd'T'HH:mm:ss.SSS'Z'Z".

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `tasks` - The secmaster tasks list.

  The [tasks](#tasks_struct) structure is documented below.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `id` - The task ID.

* `aopengine_task_id` - The engine task ID.

* `name` - The task name.

* `project_id` - The project ID.

* `description` - The task description.

* `create_time` - The creation time.

* `creator_id` - The creator ID.

* `creator_name` - The creator name.

* `update_time` - The update time.

* `modifier_id` - The modifier ID.

* `modifier_name` - The modifier name.

* `approveuser_id` - The approver user ID.

* `approveuser_name` - The approver user name.

* `approver` - The approver name.

* `notes` - The task notes.

* `definition_key` - The node workflow topology key.

* `note` - The task note.

* `due_date` - The due date.

* `action_id` - The workflow or playbook ID.

* `action_version_id` - The workflow or playbook version ID.

* `action_instance_id` - The workflow or playbook instance ID.

* `workspace_id` - The workspace ID.

* `review_comments` - The review comments.

* `view_parameters` - The view parameters.

* `handle_parameters` - The handle parameters.

* `business_type` - The business type.

* `related_object` - The related workflow or playbook name.

* `attachment_id_list` - The attachment ID list.

* `comments` - The task comments.

  The [comments](#tasks_comments_struct) structure is documented below.

* `status` - The task status. Supported values are **waiting**, **canceled** and **finished**.

* `due_handle` - The due handle method. Supported values are **CONTINUE** and **INTERRUPT**.

<a name="tasks_comments_struct"></a>
The `comments` block supports:

* `id` - The comment ID.

* `message` - The comment message.

* `user_id` - The comment user ID.

* `user_name` - The comment user name.

* `time` - The comment time.
