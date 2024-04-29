---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_organizational_policy_assignments"
description: ""
---

# huaweicloud_rms_organizational_policy_assignments

Use this data source to get the list of RMS organizational policy assignments.

## Example Usage

```hcl
variable "policy_assignment_name" {}

data "huaweicloud_organizations_organization" "test" {}

data "huaweicloud_rms_organizational_policy_assignments" "test" {
  organization_id = data.huaweicloud_organizations_organization.test.id
  name            = var.policy_assignment_name
}
```

## Argument Reference

The following arguments are supported:

* `organization_id` - (Required, String) Specifies the ID of the organization.

* `name` - (Optional, String) Specifies the name of the organizational policy assignment.

* `assignment_id` - (Optional, String) Specifies the ID of the organizational policy assignment.

* `policy_definition_id` - (Optional, String) Specifies the ID of the policy definition.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `assignments` - The organization assignments.

  The [assignments](#assignments_struct) structure is documented below.

<a name="assignments_struct"></a>
The `assignments` block supports:

* `organization_id` - The ID of the organization.

* `owner_id` - The creator of the organizational policy assignment.

* `name` - The name of the organizational policy assignment.

* `id` - The ID of the organizational policy assignment.

* `description` - The description of the organizational policy assignment.

* `policy_definition_id` - The ID of the built-in policy definition.

* `period` - The trigger period of the organizational policy assignment.

* `organization_policy_assignment_urn` - The unique identifier of the organizational policy assignment.

* `policy_filter` - The configuration used to filter resources.

  The [policy_filter](#assignments_policy_filter_struct) structure is documented below.

* `parameters` - The rule definition of the organizational policy assignment.

* `created_at` - The creation time of the organizational policy assignment.

* `updated_at` - The latest update time of the organizational policy assignment.

<a name="assignments_policy_filter_struct"></a>
The `policy_filter` block supports:

* `region_id` - The name of the region to which the filtered resources belong.

* `resource_provider` - The service name to which the filtered resources belong.

* `resource_type` - The resource type of the filtered resources.

* `resource_id` - The resource ID used to filter a specified resource.

* `tag_key` - The tag name used to filter resources.

* `tag_value` - The tag value used to filter resources.
