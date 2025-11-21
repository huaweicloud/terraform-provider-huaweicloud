---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_drift_details"
description: |-
  Use this data source to list the drift information in Resource Governance Center.
---

# huaweicloud_rgc_drift_details

Use this data source to list the drift information in Resource Governance Center.

## Example Usage

```hcl
data "huaweicloud_rgc_drift_details" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attributes Reference

* `id` - The data source ID.

* `drift_details` - Information about the detail information of drift detect.

The [drift_details](#drift_details) structure is documented below.

<a name="drift_details"></a>
The `drift_details` block supports:

* `drift_message` - The drift message describing the drift.

* `drift_target_id` - The ID of the account or OU where the drift occurred.
  
* `drift_target_type` - The type of drift target, which can be `account`, `ou`, or `policy`.
  
* `drift_type` - The type of drift.
  
* `managed_account_id` - The ID of the managed account.
  
* `parent_organizational_unit_id` - The ID of the parent organizational unit.

* `solve_solution` - The solution to resolve the drift.
