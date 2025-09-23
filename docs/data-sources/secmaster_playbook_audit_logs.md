---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook_audit_logs"
description: |-
  Use this data source to get the list of SecMaster playbook audit logs.
---

# huaweicloud_secmaster_playbook_audit_logs

Use this data source to get the list of SecMaster playbook audit logs.

## Example Usage

```hcl
variable "workspace_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_secmaster_playbook_audit_logs" "test" {
  workspace_id = var.workspace_id
  start_time   = var.start_time
  end_time     = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `start_time` - (Optional, String) Specifies the start time.
  For example, **2024-09-26T15:04:05.000Z+0800**.

* `end_time` - (Optional, String) Specifies the end time.
  For example, **2024-09-26T15:04:05.000Z+0800**.

* `status` - (Optional, String) Specifies the status.
  The value can be **RUNNING**, **FINISHED**, **FAILED**, **RETRYING**, or **TERMINATED**.

* `instance_type` - (Optional, String) Specifies the instance type. The value can be **AOP_WORKFLOW**, **SCRIPT**, or **PLAYBOOK**.

* `action_name` - (Optional, String) Specifies the workflow name.

* `instance_id` - (Optional, String) Specifies the instance ID.

* `log_level` - (Optional, String) Specifies the log level. The value can be **DEBUG**, **INFO**, **WARN** or **ERROR**.

* `trigger_type` - (Optional, String) Specifies the triggering type. The valid values are as follows:
  + **TIMER**: indicates scheduled triggering.
  + **EVENT**: indicates event triggering.

* `action_id` - (Optional, String) Specifies the workflow ID.

* `parent_instance_id` - (Optional, String) Specifies the instance ID of the parent node.

* `input` - (Optional, String) Specifies the input information.

* `output` - (Optional, String) Specifies the output information.

* `error_msg` - (Optional, String) Specifies the error message.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `audit_logs` - The audit log list.

  The [audit_logs](#audit_logs_struct) structure is documented below.

<a name="audit_logs_struct"></a>
The `audit_logs` block supports:

* `status` - The status.

* `instance_type` - The instance type.

* `action_name` - The workflow name.

* `instance_id` - The instance ID.

* `log_level` - The log level.

* `start_time` - The start time.

* `end_time` - The end time.

* `trigger_type` - The triggering type.

* `action_id` - The workflow ID.

* `parent_instance_id` - The instance ID of the parent node.

* `input` - The input information.

* `output` - The output information.

* `error_msg` - The error message.
