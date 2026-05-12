---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eip_bandwidth_associate"
description: |-
  Manages an EIP bandwidth associate resource to attach or detach an EIP from a shared bandwidth in HuaweiCloud.
---

# huaweicloud_eip_bandwidth_associate

Manages an EIP bandwidth associate resource to attach or detach an EIP from a shared bandwidth in HuaweiCloud.

## Example Usage

```hcl
variable "bandwidth_id" {} 
variable "eip_id" {}
variable "bandwidth_charge_mode" {}
variable "bandwidth_size" {}

resource "huaweicloud_eip_bandwidth_associate" "test" { 
  publicip_id  = var.eip_id 
  bandwidth_id = var.bandwidth_id

  bandwidth_charge_mode = var.bandwidth_charge_mode
  bandwidth_size        = var.bandwidth_size
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the VPC EIP associate resource. If
  omitted, the provider-level region will be used. Changing this creates a new resource.

* `publicip_id` - (Required, String, NonUpdatable) Specifies the ID of the EIP to be associated with the shared bandwidth.

* `bandwidth_id` - (Required, String, NonUpdatable) Specifies the ID of the shared bandwidth.

* `bandwidth_charge_mode` - (Required, String, NonUpdatable) Specifies the charging mode for the bandwidth after the
  EIP is removed from the shared bandwidth.

* `bandwidth_size` - (Required, Int, NonUpdatable) Specifies the bandwidth size (Mbit/s) after the EIP is removed from the
  shared bandwidth.

* `bandwidth_name` - (Optional, String, NonUpdatable) Specifies the bandwidth name after the EIP is removed from the shared
  bandwidth.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, the format is`<publicip_id>/<bandwidth_id>`.

* `public_ip_address` - The IPv4 address of the EIP.

* `public_ipv6_address` - The IPv6 address of the EIP.

* `publicip_type` - The type of the EIP.

* `ip_version` - The IP version of the EIP.

## Import

The EIP bandwidth associate can be imported using the `publicip_id` and `bandwidth_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_eip_bandwidth_associate.test <publicip_id>/<bandwidth_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include `bandwidth_charge_mode`, `bandwidth_size` and `bandwidth_name`. It is generally
recommended running `terraform plan` after importing a resource. You can then decide if changes should be applied to the
resource, or the resource definition should be updated to align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_eip_bandwidth_associate" "test" {
  ...

  lifecycle {
    ignore_changes = [
      bandwidth_charge_mode,
      bandwidth_size,
      bandwidth_name,
    ]
  }
}
```
