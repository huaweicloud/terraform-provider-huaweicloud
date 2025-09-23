---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_aggregator_policy_states"
description: |-
  Use this data source to get the list of RMS resource aggregator policy states.
---

# huaweicloud_rms_resource_aggregator_policy_states

Use this data source to get the list of RMS resource aggregator policy states.

## Example Usage

```hcl
variable "aggregator_id" {}

data "huaweicloud_rms_resource_aggregator_policy_states" "test" {
  aggregator_id = var.aggregator_id
}
```

## Argument Reference

The following arguments are supported:

* `aggregator_id` - (Required, String) Specifies the aggregator ID.

* `account_id` - (Optional, String) Specifies the ID of account to which the resource belongs.

* `policy_assignment_name` - (Optional, String) Specifies the policy assignment name.

* `compliance_state` - (Optional, String) Specifies the compliance state.
  The value can be: **Compliant** and **NonCompliant**.

* `resource_name` - (Optional, String) Specifies the resource name.

* `resource_id` - (Optional, String) Specifies the resource ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `states` - The policy states list.

  The [states](#states) structure is documented below.

<a name="states"></a>
The `states` block supports:

* `domain_id` - The domain ID.

* `region_id` - The ID of the region the resource belongs to.

* `resource_id` - The resource ID.

* `resource_name` - The resource name.

* `resource_provider` - The cloud service name.

* `resource_type` - The resource type.

* `trigger_type` - The trigger type. The value can be **resource** or **period**.

* `compliance_state` - The compliance status.

* `policy_assignment_id` - The policy assignment ID.

* `policy_assignment_name` - The policy assignment name.

* `policy_definition_id` - The ID of the policy definition.

* `evaluation_time` - The evaluation time of compliance status.
