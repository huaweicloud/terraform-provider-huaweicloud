---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_antivirus_create_virus_scan_task"
description: |-
  Manages an HSS virus scan task creation resource within HuaweiCloud.
---

# huaweicloud_hss_antivirus_create_virus_scan_task

Manages an HSS virus scan task creation resource within HuaweiCloud.

-> This resource is a one-time action resource using to create HSS virus scan task. Deleting this resource will not
  clear the corresponding request record, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "task_name" {}
variable "scan_type" {}
variable "action" {}
variable "host_ids" {
  type = list(string)
}

resource "huaweicloud_hss_antivirus_create_virus_scan_task" "test" {
  task_name = var.task_name
  scan_type = var.scan_type
  action    = var.action
  host_ids  = var.host_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `task_name` - (Required, String, NonUpdatable) Specifies the task name.

* `scan_type` - (Required, String, NonUpdatable) Specifies the task type.
  The valid values are as follows:
  + **quick**: Quick scan.
  + **full**: Full scan.
  + **custom**: Custom scan.

* `action` - (Required, String, NonUpdatable) Specifies disposal action.  
  The valid values are as follows:
  + **auto**: Automatic disposal.
  + **manual**: Manual disposal.

* `host_ids` - (Required, List, NonUpdatable) Specifies the host IDs.

* `task_id` - (Optional, String, NonUpdatable) Specifies the task ID. When creating a virus scanning task, the `task_id`
  is null, when rescanning, the `task_id` is not null, it is the ID of the current task.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
