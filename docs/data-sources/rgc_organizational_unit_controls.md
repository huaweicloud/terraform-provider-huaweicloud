---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_organizational_unit_controls"
description: |-
  Use this data source to list organizational unit controls in Resource Governance Center.
---

# huaweicloud_rgc_organizational_unit_controls

Use this data source to list organizational unit controls in Resource Governance Center.

## Example Usage

```hcl
variable managed_organizational_unit_id {}

data "huaweicloud_rgc_organizational_unit_controls" "test" {
  managed_organizational_unit_id = var.managed_organizational_unit_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `managed_organizational_unit_id` - (Required, String) The ID of the managed organizational unit.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `control_summaries` - Information about the list of control summaries.

  The [control_summaries](#control_summaries) structure is documented below.

<a name="control_summaries"></a>
The `control_summaries` block supports:

* `manage_account_id` - The ID of the manage account.

* `control_identifier` - The identifier of the control.

* `state` - The state of the control.

* `version` - The version of the control.

* `name` - The name of the control.

* `description` - The description of the control.

* `control_objective` - The objective of the control.

* `behavior` - The behavior of the control.

* `owner` - The owner of the control.

* `regional_preference` - The regional preference of the control.

* `guidance` - The guidance for the control.

* `service` - The service associated with the control.

* `implementation` - The implementation details of the control.
