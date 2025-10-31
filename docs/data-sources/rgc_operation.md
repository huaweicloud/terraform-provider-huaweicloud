---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_operation"
description: |-
 Use this data source to get operations in Resource Governance Center.
---

# huaweicloud_rgc_operation

Use this data source to get operations in Resource Governance Center.

## Example Usage

```hcl
variable organizational_unit_id {}
variable account_id {}

data "huaweicloud_rgc_operation" "test_organizational_unit_id" {
  organizational_unit_id = var.organizational_unit_id
}

data "huaweicloud_rgc_operation" "test_account" {
  account_id = var.account_id
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `account_id` - (Optional, String) The ID of the account associated with the RGC operation.

* `organizational_unit_id` - (Optional, String) The ID of the organizational unit associated with the RGC operation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `operation_id` - The ID of the RGC operation.

* `percentage_complete` - The percentage of the operation that has been completed.

* `status` - The current status of the RGC operation.

* `percentage_details` - A list of details about the percentage completion of the operation.

The [percentage_details](#percentage_details) structure is documented below.

<a name="percentage_details"></a>
The `percentage_details` block supports:

* `percentage_name` - The name of the percentage detail.

* `percentage_status` - The status of the percentage detail.
