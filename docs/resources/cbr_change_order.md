---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cbr_change_order"
description: |-
  Using this resource to change the order of a CBR resource within HuaweiCloud.
---

# huaweicloud_cbr_change_order

Using this resource to change the order of a CBR resource within HuaweiCloud.

-> This resource is only a one-time action resource to change the order of a CBR resource. Deleting this resource will
**not** revert the order change, it only removes the resource information from the TF state.

-> The resource to change order must be in **prePaid** mode. Using this resource may cause unexpected changes to
the `size` field of the `huaweicloud_cbr_vault` resource. At this time, you can use the `lifecycle` statement
to ignore the change of `size`.

## Example Usage

```hcl
variable "resource_id" {}
variable "product_id" {}
variable "resource_size" {}
variable "resource_spec_code" {}

resource "huaweicloud_cbr_change_order" "test" {
  resource_id = var.resource_id

  product_info {
    product_id               = var.product_id
    resource_size            = var.resource_size
    resource_size_measure_id = 17
    resource_spec_code       = var.resource_spec_code
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `resource_id` - (Required, String, NonUpdatable) Specifies the ID of the CBR resource whose order will be changed.
  The resource must be in **prePaid** mode.

* `product_info` - (Required, List, NonUpdatable) Specifies the product information.
  The [product_info](#change_order_product_info_struct) structure is documented below.

* `promotion_info` - (Optional, String, NonUpdatable) Specifies the promotion information for the order.

<a name="change_order_product_info_struct"></a>
The `product_info` block supports:

* `product_id` - (Required, String, NonUpdatable) Specifies the product ID, which is obtained through the price query API.
  The value consists of `1` to `64` characters and can contain only letters, digits, underscores (_), and hyphens (-).

* `resource_size` - (Required, Int, NonUpdatable) Specifies the size of the resource. Value range: `10`-`10,485,760`.
  The size value required must be greater than the existing size value.

* `resource_size_measure_id` - (Required, Int, NonUpdatable) Specifies the measurement unit of the resource size.
  Currently, only `17` (GB) is supported.

* `resource_spec_code` - (Required, String, NonUpdatable) Specifies the spec code of the resource.
  Valid values are: **vault.backup.server.normal**, **vault.backup.turbo.normal**, **vault.backup.database.normal**,
  **vault.backup.volume.normal**, **vault.backup.rds.normal**, **vault.replication.server.normal**, and
  **vault.hybrid.server.normal**.

## Attribute Reference

The following attributes are exported:

* `id` - The ID of the resource.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
