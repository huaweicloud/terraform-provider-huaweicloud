---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_enabled_controls"
description: |-
 Use this data source to list enabled controls in Resource Governance Center.
---

# huaweicloud_rgc_enabled_controls

Use this data source to list enabled controls in Resource Governance Center.

## Example Usage

```hcl
data "huaweicloud_rgc_enabled_controls" "enabled_controls" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `enabled_controls` - Information about the enabled control list.

The [enabled_controls](#enabled_controls) structure is documented below.

<a name="enabled_controls"></a>
The `enabled_controls` block supports:

* `manage_account_id` - The ID of the management account.

* `control_identifier` - The identifier of the control.

* `name` - The name of the control.

* `description` - The description of the enabled control.

* `control_objective` - The objective of the enabled control.

* `behavior` - The behavior of the enabled control.

* `owner` - The owner of the enabled control.

* `regional_preference` - The regional preference of the enabled control.
