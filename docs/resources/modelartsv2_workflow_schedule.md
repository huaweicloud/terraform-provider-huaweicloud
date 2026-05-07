---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_workflow_schedule"
description: |-
  Manages a ModelArts workflow schedule configuration within HuaweiCloud.
---

# huaweicloud_modelartsv2_workflow_schedule

Manages a ModelArts workflow schedule configuration within HuaweiCloud.

## Example Usage

```hcl
variable "workflow_id" {}

resource "huaweicloud_modelartsv2_workflow_schedule" "test" {
  workflow_id = var.workflow_id

  content = jsonencode({
    cron   = "0 0 0 * * Thu"
    method = "fixed"
  })
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the workflow schedule is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `workflow_id` - (Required, String, NonUpdatable) Specifies the ID of the workflow to which the schedule configuration
  belongs.

* `content` - (Required, String) Specifies the content of the workflow schedule, in JSON format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `enable` - Whether the workflow schedule is enabled.

* `policies` - The scheduling policies of the workflow schedule.  
  The [policies](#modelartsv2_workflow_schedule_policies_attr) structure is documented below.

* `user_id` - The user ID that created the workflow schedule.

* `created_at` - The creation time of the workflow schedule, in RFC3339 format.

<a name="modelartsv2_workflow_schedule_policies_attr"></a>
The `policies` block supports:

* `on_failure` - The policy action when the workflow execution fails.

* `on_running` - The policy action when the workflow is already running.

## Import

Workflow schedules can be imported using `workflow_id` and `id`, separated by a slash, e.g.

```bash
terraform import huaweicloud_modelartsv2_workflow_schedule.test <workflow_id>/<id>
```
