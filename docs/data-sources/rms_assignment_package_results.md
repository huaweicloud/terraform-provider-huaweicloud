---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_assignment_package_results"
description: |-
  Use this data source to get the list of RMS assignment package results.
---

# huaweicloud_rms_assignment_package_results

Use this data source to get the list of RMS assignment package results.

## Example Usage

```hcl
variable "assignment_package_id" {}

data "huaweicloud_rms_assignment_package_results" "basic" {
  assignment_package_id = var.assignment_package_id
}
```

## Argument Reference

The following arguments are supported:

* `assignment_package_id` - (Required, String) Specifies the assignment package name.

* `policy_assignment_name` - (Optional, String) Specifies the policy assignment name. Fuzzy search is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `value` - The details about compliance results of assignments in an assignment package.

  The [value](#value_struct) structure is documented below.

<a name="value_struct"></a>
The `value` block supports:

* `policy_assignment_id` - The policy assignment ID.

* `policy_assignment_name` - The policy assignment name.

* `resource_id` - The ID of the resource to be evaluated.

* `compliance_state` - The compliance result of the assignment.

* `evaluation_time` - The time for evaluating resources.
