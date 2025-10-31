---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_blueprint"
description: |-
  Use this data source to get blue print in Resource Governance Center.
---

# huaweicloud_rgc_blueprint

Use this data source to get blue print in Resource Governance Center.

## Example Usage

```hcl
variable "managed_account_id" {}

data "huaweicloud_rgc_blueprint" "test" {
  managed_account_id = var.managed_account_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `managed_account_id` - (Required) The ID of the managed account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `manage_account_id` - The ID of the management account.

* `account_id` - The ID of the managed account.

* `account_name` - The name of the managed account.

* `blueprint_product_id` - The ID of the blueprint product.

* `blueprint_product_name` - The name of the blueprint product.

* `blueprint_product_version` - The version of the blueprint product.

* `blueprint_status` - The status of the blueprint deployment, including **failed**, **completed**, and **in progress**.

* `blueprint_product_impl_type` - The implementation type of the blueprint product.

* `blueprint_product_impl_detail` - The implementation details of the blueprint product.
