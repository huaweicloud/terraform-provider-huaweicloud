---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_policy_assignment_envaluate"
description: |-
  Manages a RMS policy assignment evaluate resource within HuaweiCloud resources.
---

# huaweicloud_rms_policy_assignment_evaluate

Manages a RMS policy assignment evaluate resource within HuaweiCloud resources.

## Example Usage

```hcl
variable "policy_assignment_id" {}

resource "huaweicloud_rms_policy_assignment_envaluate" "test" {
  policy_assignment_id = var.policy_assignment_id
}
```

## Argument Reference

The following arguments are supported:

* `policy_assignment_id` - (Required, String, ForceNew) Specifies the ID of the policy assignment to evaluate.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the policy assignment evaluate.
