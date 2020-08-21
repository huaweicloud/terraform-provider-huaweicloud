---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_eip_associate"
sidebar_current: "docs-huaweicloud-resource-compute-eip-associate"
description: |-
  Associate a EIP to an instance
---

# huaweicloud\_compute\_eip_associate

Associate an EIP to an instance. This is an alternative to
`huaweicloud_compute_floatingip_associate_v2`.

## Example Usage

### Automatically detect the correct network

```hcl
resource "huaweicloud_compute_instance" "instance_1" {
  name            = "instance_1"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = 3
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]
}

resource "huaweicloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_compute_eip_associate" "associate_1" {
  public_ip   = huaweicloud_vpc_eip_v1.eip_1.address
  instance_id = huaweicloud_compute_instance_v2.instance_1.id
}
```

### Explicitly set the network to attach to

```hcl
resource "huaweicloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = 3
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    name = "my_network"
  }

  network {
    name = "default"
  }
}

resource "huaweicloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "test"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_compute_eip_associate" "fip_1" {
  public_ip   = huaweicloud_vpc_eip_v1.eip_1.address
  instance_id = huaweicloud_compute_instance_v2.instance_1.id
  fixed_ip    = huaweicloud_compute_instance_v2.instance_1.network.1.fixed_ip_v4
}
```

## Argument Reference

The following arguments are supported:

* `public_ip` - (Required) The EIP to associate.

* `floating_ip` - (Deprecated) Use `public_ip` instead. The EIP to associate.

* `instance_id` - (Required) The instance to associte the EIP with.

* `fixed_ip` - (Optional) The specific IP address to direct traffic to.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `floating_ip` - Deprecated. See Argument Reference above.
* `public_ip` - See Argument Reference above.
* `instance_id` - See Argument Reference above.
* `fixed_ip` - See Argument Reference above.

## Import

This resource can be imported by specifying all three arguments, separated
by a forward slash:

```
$ terraform import huaweicloud_compute_eip_associate.eip_1 <eip>/<instance_id>/<fixed_ip>
```
