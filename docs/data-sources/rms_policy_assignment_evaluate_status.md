---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_policy_assignment_evaluate_status"
description: |-
  Use this data source to get the evaluation status of a policy assignment.
---

# huaweicloud_rms_policy_assignment_evaluate_status

Use this data source to get the evaluation status of a policy assignment.

## Example Usage

```hcl
variable "policy_assignment_id" {}

data "huaweicloud_rms_policy_assignment_evaluate_status" "test" {
  policy_assignment_id = var.policy_assignment_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_assignment_id` - (Required, String) Specifies the policy assignment ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `start_time` - Indicates the start time of the evaluation task.

* `end_time` - Indicates the end time of the evaluation task.

* `state` - Indicates the execution status of the evaluation task.

* `error_message` - Indicates the failure information of the evaluation task.
