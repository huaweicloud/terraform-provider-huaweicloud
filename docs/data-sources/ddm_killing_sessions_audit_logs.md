---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_killing_sessions_audit_logs"
description: |-
  Use this data source to get the list of killing sessions audit log.
---

# huaweicloud_ddm_killing_sessions_audit_logs

Use this data source to get the list of killing sessions audit log.

## Example Usage

```hcl
variable "instance_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_ddm_killing_sessions_audit_logs" "test" {
  instance_id = var.instance_id
  start_time  = var.start_time
  end_time    = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the DDM instance ID or ID of the associated RDS instance.

* `start_time` - (Required, String) Specifies the start time in UTC, accurate to milliseconds.
  The format is **yyyy-mm-ddThh:mm:ssZ**.

* `end_time` - (Required, String) Specifies the end time in UTC, accurate to milliseconds.
  The format is **yyyy-mm-ddThh:mm:ssZ**. The interval between the start time and the end time must be no more than 7 days.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `process_audit_logs` - Indicates the list of killing sessions audit log.

  The [process_audit_logs](#process_audit_logs_struct) structure is documented below.

<a name="process_audit_logs_struct"></a>
The `process_audit_logs` block supports:

* `process_id` - Indicates the session ID.

* `instance_id` - Indicates the DDM instance ID.

* `instance_name` - Indicates the DDM instance name.

* `execute_time` - Indicates the execute time in UTC format.

* `execute_user_name` - Indicates the name of the user who executes the task.
