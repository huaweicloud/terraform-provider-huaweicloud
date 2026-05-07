---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_workflow_schedules"
description: |-
  Use this data source to get the list of workflow schedules within HuaweiCloud.
---

# huaweicloud_modelartsv2_workflow_schedules

Use this data source to get the list of workflow schedules within HuaweiCloud.

## Example Usage

```hcl
variable "workflow_id" {}

data "huaweicloud_modelartsv2_workflow_schedules" "test" {
  workflow_id = var.workflow_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the workflow schedules are located.  
  If omitted, the provider-level region will be used.

* `workflow_id` - (Required, String) Specifies the ID of the workflow to which the schedule configurations belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `schedules` - The list of the workflow schedules.  
  The [schedules](#modelartsv2_workflow_schedules_attr) structure is documented below.

<a name="modelartsv2_workflow_schedules_attr"></a>
The `schedules` block supports:

* `id` - The ID of the workflow schedule.

* `type` - The type of the workflow schedule.

* `content` - The content of the workflow schedule, in JSON format.

* `action` - The action of the workflow schedule.

* `workflow_id` - The ID of the workflow to which the schedule belongs.

* `user_id` - The user ID that created the workflow schedule.

* `enable` - Whether the workflow schedule is enabled.

* `policies` - The scheduling policies of the workflow schedule.  
  The [policies](#modelartsv2_workflow_schedules_policies_attr) structure is documented below.

* `created_at` - The creation time of the workflow schedule, in RFC3339 format.

<a name="modelartsv2_workflow_schedules_policies_attr"></a>
The `policies` block supports:

* `on_failure` - The policy action when the workflow execution fails.

* `on_running` - The policy action when the workflow is already running.
