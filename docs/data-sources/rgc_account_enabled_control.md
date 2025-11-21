---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_account_enabled_control"
description: |-
  Use this data source to get enabled control information for an managed account in Resource Governance Center.
---

# huaweicloud_rgc_account_enabled_control

Use this data source to get enabled control information for an managed account in Resource Governance Center.

## Example Usage

```hcl
variable "control_id" {}
variable "managed_account_id" {}

data "huaweicloud_rgc_account_enabled_control" "test" {
  control_id         = var.control_id
  managed_account_id = var.managed_account_id
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `control_id` - (Required, String) The ID of the control policy.

* `managed_account_id` - (Required, String) The ID of the managed account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `control_detail` - Information about the enabled control policy.
  The [control_detail](#control_detail_struct) structure is documented below.

* `regions` - Information about the enabled control policy.
  The [regions](#regions_struct) structure is documented below.

* `state` - The enablement status of the control policy.

* `message` - The message of the control policy.

* `version` - The current version number of the control policy.

<a name="control_detail_struct"></a>
The `control_detail` block supports:

* `manage_account_id` - The ID of the management account.

* `control_identifier` - The identifier of the control.

* `name` - The name of the control policy.

* `description` - The description of the control policy.

* `control_objective` - The objective of the control policy.

* `behavior` - The type of the control policy. It can be **Proactive**, **Detective** or **Preventive**.

* `owner` - The owner of the control policy.

* `regional_preference` - The regional preference of the control.

<a name="regions_struct"></a>
The `regions` block supports:

* `region` - The name of the region.

* `region_configuration_status` - The status of the region.
