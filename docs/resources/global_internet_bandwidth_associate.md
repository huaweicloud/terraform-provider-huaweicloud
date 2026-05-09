---
subcategory: "EIP (Elastic IP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_global_eip_internet_bandwidth_associate"
description: |-
  Manages a Global EIP associate internet bandwidth resource within HuaweiCloud.
---

# huaweicloud_global_internet_bandwidth_associate

Manages a Global EIP associate internet bandwidth resource within HuaweiCloud.

## Example Usage

```hcl
variable "global_eip_id" {}
variable "internet_bandwidth_id" {}

resource "huaweicloud_global_internet_bandwidth_associate" "test" {
  global_eip_id = var.global_eip_id
  
  global_eip {
    internet_bandwidth_id = var.internet_bandwidth_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the resource is located.
  If omitted, the provider-level region will be used. Change this parameter will create a new resource.

* `global_eip_id` - (Required, String, NonUpdatable) Specifies the ID of the global EIP to which the internet bandwidth
  will be associated.

* `global_eip` - (Required, List, NonUpdatable) Specifies the global EIP configuration for associating internet
  bandwidth. The object structure is documented below.

  The `global_eip` block supports:

  * `internet_bandwidth_id` - (Required, String, NonUpdatable) Specifies the ID of the internet bandwidth to associate
    to the global EIP.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in the format of `<global_eip_id>`.

* `global_eip/internet_bandwidth_info` - The information about the associated internet bandwidth. The object structure
  is documented below.

  The `internet_bandwidth_info` block exports:

  * `id` - The ID of the attached internet bandwidth.

  * `size` - The size of the attached internet bandwidth.

## Import

This resource can be imported using the `global_eip_id`, e.g.

```bash
$ terraform import huaweicloud_global_internet_bandwidth_associate.test <global_eip_id>
```
