---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_organizational_units"
description: |-
 Use this data source to list organizational units in Resource Governance Center.
---

# huaweicloud_rgc_organizational_units

Use this data source to list organizational units in Resource Governance Center.

## Example Usage

```hcl
data "huaweicloud_rgc_organizational_units" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `control_id` - (Optional, String) The ID of the enabled control policy. It can be used to filter out the
  organizational units that have this control enabled. It must be a string of 1 to 128 characters,
  which can include letters, numbers, underscores, and hyphens.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `managed_organizational_units` - Information about the managed organizational units list.

The [managed_organizational_units](#managed_organizational_units) structure is documented below.

<a name="managed_organizational_units"></a>
The `managed_organizational_units` block supports:

* `landing_zone_version` - The version of the Landing Zone.

* `organizational_unit_status` - The status of the organizational unit.

* `organizational_unit_type` - The type of the organizational unit. Possible values are `CORE`, `CUSTOM`, and `ROOT`.

* `manage_account_id` - The ID of the management account.

* `organizational_unit_id` - The ID of the organizational unit.

* `organizational_unit_name` - The name of the organizational unit.

* `parent_organizational_unit_id` - The ID of the parent organizational unit.

* `parent_organizational_unit_name` - The name of the parent organizational unit.

* `created_at` - The time when the organizational unit was created.

* `message` - The error message description.

* `updated_at` - The time when the organizational unit was last updated.
