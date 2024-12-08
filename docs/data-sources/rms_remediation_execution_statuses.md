---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_remediation_execution_statuses"
description: |-
  Use this data source to get the list of RMS latest remediation execuation statuses.
---

# huaweicloud_rms_remediation_execution_statuses

Use this data source to get the list of RMS latest remediation execution statuses.

## Example Usage

```hcl
variable "policy_assignment_id" {}

data "huaweicloud_rms_remediation_execution_statuses" "test" {
  policy_assignment_id = var.policy_assignment_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `policy_assignment_id` - (Required, String) The policy assignment ID.

* `resource_keys` - (Optional, List) The list of query criteria required to collect remediation results.

  The [resource_keys](#resource_keys_struct) structure is documented below.

<a name="resource_keys_struct"></a>
The `resource_keys` block supports:

* `resource_type` - (Required, String) The resource type.

* `resource_id` - (Required, String) The resource ID.

* `resource_provider` - (Required, String) The cloud service name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `value` - The compliance rule remediation execution results.

  The [value](#value_struct) structure is documented below.

<a name="value_struct"></a>
The `value` block supports:

* `resource_key` - The query criteria required to collect remediation results.

  The [resource_key](#value_resource_key_struct) structure is documented below.

* `invocation_time` - The start time of remediation.

* `state` - The execution state of remediation.

* `message` - The information of remediation execution.

<a name="value_resource_key_struct"></a>
The `resource_key` block supports:

* `resource_type` - The resource type.

* `resource_id` - The resource ID.

* `resource_provider` - The cloud service name.
