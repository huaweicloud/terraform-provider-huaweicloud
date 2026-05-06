---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_eip_update_publicip_pool"
description: |-
  Manages a VPC EIP public IP pool resource within HuaweiCloud.
---

# huaweicloud_vpc_eip_update_publicip_pool

Manages a VPC EIP public IP pool resource within HuaweiCloud.

-> This resource is used to update a public IP pool.
Deleting this resource will not delete the actual public IP pool on the cloud, but will only remove the resource
information from the tfstate file.

## Example Usage

```hcl
variable "publicip_pool_id" {}
variable "name" {}
variable "description" {}

resource "huaweicloud_vpc_eip_update_publicip_pool" "test" {
  publicip_pool_id = var.publicip_pool_id
  name             = var.name
  description      = var.description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the VPC EIP public IP pool resource. If
  omitted, the provider-level region will be used. Changing this creates a new resource.

* `publicip_pool_id` - (Required, String, NonUpdatable) Specifies the ID of the public IP pool to be updated.

* `name` - (Optional, String) Specifies the name of the public IP pool. The name needs to be unique.

* `description` - (Optional, String) Specifies the description of the public IP pool.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the same as the `publicip_pool_id`.

* `status` - The status of the public IP pool. Valid values are **ACTIVE**, **INACTIVE**, and **SOLDOUT**.

* `type` - The type of the public IP pool. Valid values are **spec_bgp** and **spec_sbgp**.

* `project_id` - The project ID of the public IP pool.

* `size` - The size of the public IP pool.

* `used` - The number of used IPs in the public IP pool.

* `created_at` - The creation time of the public IP pool.

* `updated_at` - The last update time of the public IP pool.

* `billing_info` - The billing information of the public IP pool.

  The [billing_info](#billing_info_struct) structure is documented below.

* `public_border_group` - The public border group of the public IP pool.

* `shared` - Whether the public IP pool is shared.

* `tags` - The tags associated with the public IP pool.

  The [tags](#tags_struct) structure is documented below.

* `enterprise_project_id` - The enterprise project ID of the public IP pool.

* `allow_share_bandwidth_types` - The allowed shared bandwidth types for the public IP pool.

<a name="billing_info_struct"></a>
The `billing_info` block supports:

* `order_id` - The order ID associated with the public IP pool.

* `product_id` - The product ID associated with the public IP pool.

<a name="tags_struct"></a>
The `tags` block supports:

* `key` - The tag key.

* `value` - The tag value.

## Import

The public IP pool can be imported using the `publicip_pool_id`:

```bash
terraform import huaweicloud_vpc_eip_update_publicip_pool.test <publicip_pool_id>
```
