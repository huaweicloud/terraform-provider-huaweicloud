---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_remediation_execution"
description: |-
  Manages a RMS remediation execution resource within HuaweiCloud.
---

# huaweicloud_rms_remediation_execution

Manages a RMS remediation execution resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not change the current status.

## Example Usage

```hcl
variable policy_assignment_id {}

resource "huaweicloud_rms_remediation_execution" "test" {
  policy_assignment_id = var.policy_assignment_id
  all_supported        = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `policy_assignment_id` - (Required, String, NonUpdatable) Specifies the policy assignment ID.

* `all_supported` - (Required, Bool, NonUpdatable) Specifies whether to perform remediation for all non-compliant resources.

* `resource_ids` - (Optional, List, NonUpdatable) Specifies the list of resource IDs that require remediation.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `result` - The result of the remediation execution.

  The [result](#result_struct) structure is documented below.

<a name="result_struct"></a>
The `result` block supports:

* `invocation_time` - The start time of remediation.

* `state` - The execution state of remediation.

* `message` - The information of remediation execution.

* `automatic` - Whether the remediation is automatic.

* `resource_id` - The resource ID.

* `resource_name` - The resource name.

* `resource_provider` - The cloud service name.

* `resource_type` - The resource type.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
