---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_policy_assignments"
description: ""
---

# huaweicloud_rms_policy_assignments

Use this data source to get the list of RMS policy assignments.

## Example Usage

```hcl
variable "policy_assignment_name" {}

data "huaweicloud_rms_policy_assignments" "test" {
  name   = var.policy_assignment_name
  status = "Enabled"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the assignment name.

* `assignment_id` - (Optional, String) Specifies the ID of the policy assignment.

* `policy_definition_id` - (Optional, String) Specifies the ID of the policy definition.

* `status` - (Optional, String) Specifies the expect status of the policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `assignments` - The policy assignment list.

  The [assignments](#assignments_struct) structure is documented below.

<a name="assignments_struct"></a>
The `assignments` block supports:

* `id` - The policy assignment ID.

* `name` - The policy assignment name.

* `description` - The policy assignment description.

* `policy_assignment_type` - Specifies the policy assignment type.

* `policy_definition_id` - The ID of the policy used by the policy assignment.

* `period` - The policy assignment period.

* `policy_filter` - The configuration used to filter resources.

  The [policy_filter](#assignments_policy_filter_struct) structure is documented below.

* `custom_policy` - The configuration of the custom policy.

  The [custom_policy](#assignments_custom_policy_struct) structure is documented below.

* `parameters` - The policy assignment parameter.

* `status` - The policy assignment status.

* `created_by` - The policy assignment creator.

* `created_at` - The creation time of the policy assignment.

* `updated` - The latest update time of the policy assignment.

<a name="assignments_policy_filter_struct"></a>
The `policy_filter` block supports:

* `region` - The name of the region to which the filtered resources belong.

* `resource_provider` - The service name to which the filtered resources belong.

* `resource_type` - The resource type of the filtered resources.

* `resource_id` - The resource ID used to filter a specified resource.

* `tag_key` - The tag name used to filter resources.

* `tag_value` - The tag value used to filter resources.

<a name="assignments_custom_policy_struct"></a>
The `custom_policy` block supports:

* `function_urn` - The function URN used to create the custom policy.

* `auth_type` - The authorization type of the custom policy.

* `auth_value` - The authorization value of the custom policy.
