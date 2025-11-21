---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_control_violations"
description: |-
  Use this data source to query the details of a control violation in Resource Governance Center.
---

# huaweicloud_rgc_control_violations

Use this data source to query the details of a control violation in Resource Governance Center.

## Example Usage

```hcl
variable account_id {}
variable organizational_unit_id {}

data "huaweicloud_rgc_control_violations" "test_account" {
  account_id = var.account_id
}

data "huaweicloud_rgc_control_violations" "test_account" {
  organizational_unit_id = var.organizational_unit_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `account_id` - (Optional, String) The managed account iD。

* `organizational_unit_id` - (Optional, String) The registered organizational unit ID。

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `control_violations` - A list of control violations.
The [control_violations](#control_violations) structure is documented below.

<a name="control_violations"></a>
The `control_violations` block supports:

* `account_id` - The compliance status of the control violation.

* `account_name` - The name of the account.

* `display_name` - The display name of control strategy.

* `name` - The name of control strategy.

* `control_id` - The ID of control strategy.

* `parent_organizational_unit_id` - The ID of the parent organizational unit associated with the non-compliant control
  strategy.

* `parent_organizational_unit_name` - The name of the parent organizational unit associated with the non-compliant
  control strategy.

* `region` - The region of the non-compliant control strategy.

* `resource` - The resource of the non-compliant control strategy.

* `resource_name` - The resource name of the non-compliant control strategy.

* `resource_type` - The resource type of the non-compliant control strategy.

* `service` - The service associated with non-compliant control strategy.
