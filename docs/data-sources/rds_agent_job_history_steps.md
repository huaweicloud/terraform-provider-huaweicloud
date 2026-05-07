---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_agent_job_history_steps"
description: |-
  Use this data source to query the agent job history steps of RDS within HuaweiCloud.
---

# huaweicloud_rds_agent_job_history_steps

Use this data source to query the agent job history steps of RDS within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "history_id" {}

data "huaweicloud_rds_agent_job_history_steps" "test" {
  instance_id = var.instance_id
  history_id  = var.history_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `history_id` - (Required, String) Specifies the ID of the agent job execution history.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `steps` - The list of execution history steps.
  The [steps](#agent_job_history_steps_attr) structure is documented below.

<a name="agent_job_history_steps_attr"></a>
The `steps` block supports:

* `step_id` - The step ID.

* `step_name` - The step name.

* `run_status` - The execution status of the job.
  The valid values are as follows:
  + **failed**
  + **succeeded**
  + **retrying**
  + **canceled**
  + **in_progress**

* `run_time` - The latest execution time, in the format **yyyy-mm-ddThh:mm:ssZ**.

* `run_duration` - The job execution duration, in the format **hh:mm:ss**.

* `message` - The execution message.
