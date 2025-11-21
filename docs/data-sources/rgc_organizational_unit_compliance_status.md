---
subcategory: "RGC"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_organizational_unit_compliance_status"
description: |
  Use this data source to get the organizational unit compliance status in Resource Governance Center.
---

# huaweicloud_rgc_organizational_unit_compliance_status

Use this data source to get the organizational unit compliance status in Resource Governance Center.

## Example Usage

```hcl
variable managed_organizational_unit_id {}

data "huaweicloud_rgc_organizational_unit_compliance_status" "test" {
  managed_organizational_unit_id = var.managed_organizational_unit_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `managed_organizational_unit_id` - (Required, String) The ID of the managed organizational unit for which to retrieve
  compliance status.

* `control_id` - (Optional, String) The ID of the control to filter the compliance status.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `compliance_status` - The compliance status of the organizational unit, including **Compliant**, **NonCompliant**.
