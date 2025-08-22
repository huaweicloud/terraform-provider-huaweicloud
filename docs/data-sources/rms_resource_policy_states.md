---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_policy_states"
description: |-
  Use this data source to query policy states of a resource.
---

# huaweicloud_rms_resource_policy_states

Use this data source to query policy states of a resource.

## Example Usage

```hcl
variable "resource_id" {}

data "huaweicloud_rms_resource_policy_states" "test" {
  resource_id = var.resource_id
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, String) Specifies the resource ID.

* `compliance_state` - (Optional, String) Specifies the compliance status.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `value` - Indicates the return value of querying the compliance result.

  The [value](#value_struct) structure is documented below.

<a name="value_struct"></a>
The `value` block supports:

* `domain_id` - Indicates the user ID.

* `region_id` - Indicates the ID of the region the resource belongs to.

* `resource_id` - Indicates the resource ID.

* `resource_name` - Indicates the resource name.

* `resource_provider` - Indicates the cloud service name.

* `resource_type` - Indicates the resource type.

* `trigger_type` - Indicates the trigger type. The value can be **resource** or **period**.

* `compliance_state` - Indicates the compliance state.

* `policy_assignment_id` - Indicates the rule ID.

* `policy_assignment_name` - Indicates the rule name.

* `policy_definition_id` - Indicates the policy ID.

* `evaluation_time` - Indicates the compliance state evaluation time.
