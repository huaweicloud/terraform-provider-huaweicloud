---
subcategory: "rms"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_policy_assignment_summary"
description: |-
  Use this data source to get the summary of a policy assignment.
---

# huaweicloud_rms_policy_assignment_summary

Use this data source to get the summary of a policy assignment.

## Example Usage

```hcl
variable "policy_assignment_id" {}

data "huaweicloud_rms_policy_assignment_summary" "test" {
  policy_assignment_id = var.policy_assignment_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_assignment_id` - (Required, String) Specifies the policy assignment ID.

* `resource_name` - (Optional, String) Specifies the resource name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `compliance_state` - Indicates the policy assignment status.

* `results` - Indicates the results of compliance summaries.
  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `resource_details` - Indicates the resource compliance summary details.
  The [resource_details](#results_details_struct) structure is documented below.

* `assignment_details` - Indicates the compliance summary details.
  The [assignment_details](#results_details_struct) structure is documented below.

<a name="results_details_struct"></a>
The `resource_details` and `assignment_details` block supports:

* `compliant_count` - Indicates the number of compliant resources.

* `non_compliant_count` - Indicates the number of non-compliant resources.
