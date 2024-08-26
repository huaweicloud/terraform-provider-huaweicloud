---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_aggregator_policy_assignments"
description: |-
  Use this data source to get the list of RMS resource aggregator policy assignments.
---

# huaweicloud_rms_resource_aggregator_policy_assignments

Use this data source to get the list of RMS resource aggregator policy assignments.

## Example Usage

```hcl
variable "aggregator_id" {}

data "huaweicloud_rms_resource_aggregator_policy_assignments" "test" {
  aggregator_id = var.aggregator_id
}
```

## Argument Reference

The following arguments are supported:

* `aggregator_id` - (Required, String) Specifies the aggregator ID.

* `filter` - (Optional, List) Specifies the filter. The [filter](#filter) structure is documented below.

<a name="filter"></a>
The `filter` block supports:

* `account_id` - (Optional, String) Specifies the ID of account to which the resource belongs.

* `policy_assignment_name` - (Optional, String) Specifies the policy assignment name.

* `compliance_state` - (Optional, String) Specifies the compliance state.
  The value can be: **Compliant** and **NonCompliant**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `assignments` - The policy assignments list. The [assignments](#assignments) structure is documented below.

<a name="assignments"></a>
The `assignments` block supports:

* `policy_assignment_id` - The policy assignment ID.

* `policy_assignment_name` - The policy assignment name.

* `account_id` - The ID of account to which the resource belongs.

* `account_name` - The name of account to which the resource belongs.

* `compliance` - The compliance of the policy assignment. The [compliance](#compliance) structure is documented below.

<a name="compliance"></a>
The `compliance` block supports:

* `compliance_state` - The compliance status.

* `resource_details` - The resource details. The [resource_details](#resource_details) structure is documented below.

<a name="resource_details"></a>
The `resource_details` block supports:

* `compliant_count` - The number of compliant resources.

* `non_compliant_count` - The number of non-compliant resources.
