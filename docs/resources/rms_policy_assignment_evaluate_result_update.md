---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_policy_assignment_evaluate_result_update"
description: |-
  Manages a Config policy assignment evaluate result update resource within HuaweiCloud resources.
---

# huaweicloud_rms_policy_assignment_evaluate_result_update

Manages a Config policy assignment evaluate result update resource within HuaweiCloud resources.

## Example Usage

```hcl
variable "policy_assignment_id" {}
variable "evaluation_time" {}
variable "evaluation_hash" {}
variable "resource_id" {}
variable "resource_name" {}
variable "region_id" {}
variable "domain_id" {}

resource "huaweicloud_rms_policy_assignment_evaluate_result_update" "test" {
  policy_assignment_id = var.policy_assignment_id
  trigger_type         = "period"
  compliance_state     = "NonCompliant"
  evaluation_time      = var.evaluation_time
  evaluation_hash      = var.evaluation_hash

  policy_resource {
    resource_id       = var.resource_id
    resource_name     = var.resource_name
    resource_provider = "iam"
    resource_type     = "users"
    region_id         = var.region_id
    domain_id         = var.domain_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `policy_assignment_id` - (Required, String, NonUpdatable) Specifies the policy rule id.

* `trigger_type` - (Required, String, NonUpdatable) Specifies the trigger type.

* `compliance_state` - (Required, String, NonUpdatable) Specifies the compliance status.

* `evaluation_time` - (Required, String, NonUpdatable) Specifies when a rule is used to evaluate the resource compliance.

* `evaluation_hash` - (Required, String, NonUpdatable) Specifies the evaluation verification code.

* `policy_resource` - (Required, List, NonUpdatable) Specifies the resource.
  The [policy_resource](#policy_resource_struct) structure is documented below.

* `policy_assignment_name` - (Optional, String, NonUpdatable) Specifies the policy rule name.

<a name="policy_resource_struct"></a>
The `policy_resource` block supports:

* `resource_id` - (Optional, String, NonUpdatable) Specifies the resource ID.

* `resource_name` - (Optional, String, NonUpdatable) Specifies the resource name.

* `resource_provider` - (Optional, String, NonUpdatable) Specifies the cloud service name.

* `resource_type` - (Optional, String, NonUpdatable) Specifies the resource type.

* `region_id` - (Optional, String, NonUpdatable) Specifies the region ID.

* `domain_id` - (Optional, String, NonUpdatable) Specifies the ID of the user to which the resource belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
