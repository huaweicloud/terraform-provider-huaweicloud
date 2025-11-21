---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_organizational_unit_enabled_control"
description: |-
  Use this data source to get enabled control information for an registered organizational unit in Resource Governance Center.
---

# huaweicloud_rgc_organizational_unit_enabled_control

Use this data source to get enabled control information for an registered organizational unit in Resource Governance Center.

## Example Usage

```hcl
variable "control_id" {}
variable "managed_organizational_unit_id" {}

data "huaweicloud_rgc_organizational_unit_enabled_control" "test" {
  control_id                     = var.control_id
  managed_organizational_unit_id = var.managed_organizational_unit_id
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `control_id` - (Required, String) The ID of the control policy.

* `managed_organizational_unit_id` - (Required, String) The ID of the registered organizational unit.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `control` - Information about the enabled control policy.

* `regions` - A list of regions associated with the account.

* `state` - The enablement status of the control policy.

* `message` - The message of the control policy.

* `version` - The current version number of the control policy.

The [control](#control) structure is documented below.

<a name="control"></a>
The `control` block supports:

* `manage_account_id` - The ID of the managed account.

* `control_identifier` - The identifier of the control policy.

* `name` - The name of the control policy.

* `description` - The description information of the control policy.

* `control_objective` - The objective of the control policy.

* `behavior` - The type of control policy. It includes Proactive, Detective, and Preventive control policies.

* `owner` - The source of the managed account's creation. It includes CUSTOM and RGC.

* `regional_preference` - The regional preference, which can be either regional or global.

The [regions](#regions) structure is documented below.

<a name="regions"></a>
The `regions` block supports:

* `region` - The name of the region.

* `region_configuration_status` - The status of the region.
