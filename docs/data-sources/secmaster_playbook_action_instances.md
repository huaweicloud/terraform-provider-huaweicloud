---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook_action_instances"
description: |-
  Use this data source to get the list of SecMaster playbook action instances.
---

# huaweicloud_secmaster_playbook_action_instances

Use this data source to get the list of SecMaster playbook action instances.

## Example Usage

```hcl
variable "workspace_id" {}
variable "playbook_instance_id" {}

data "huaweicloud_secmaster_playbook_action_instances" "test" {
  workspace_id         = var.workspace_id
  playbook_instance_id = var.playbook_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID to which the playbook instance belongs.

* `playbook_instance_id` - (Required, String) Specifies the playbook instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `action_instances` - The playbook action instance list.

  The [action_instances](#action_instances_struct) structure is documented below.

<a name="action_instances_struct"></a>
The `action_instances` block supports:

* `action` - The action information of the playbook action instance.

  The [action](#action_instances_action_struct) structure is documented below.

* `instance_log` - The log information of the playbook action instance.

  The [instance_log](#action_instances_instance_log_struct) structure is documented below.

<a name="action_instances_action_struct"></a>
The `action` block supports:

* `id` - The action ID.

* `name` - The workflow name of the action.

* `description` - The description of the action.

* `action_id` - The workflow ID of the action.

* `action_type` - The workflow type of the action.

* `playbook_id` - The playbook ID associated with the action.

* `playbook_version_id` - The playbook version ID associated with the action.

<a name="action_instances_instance_log_struct"></a>
The `instance_log` block supports:

* `instance_type` - The instance type that the log printed.
  The value can be **AOP_WORKFLOW**, **SCRIPT** or **PLAYBOOK**.

* `action_name` - The workflow name that the log printed.

* `parent_instance_id` - The parent instance ID that the log printed.

* `output` - The output information that the log printed.

* `status` - The instance status that the log printed.

* `action_id` - The workflow ID that the log printed.

* `instance_id` - The instance ID that the log printed.

* `log_level` - The log level.

* `input` - The input information that the log printed.

* `error_msg` - The error message that the log printed.

* `start_time` - The start time that the log printed.

* `end_time` - The end time that the log printed.

* `trigger_type` - The triggering type that the log printed.
