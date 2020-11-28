---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud\_networking\_eip\_associate

Associates an EIP to a port. This can be used instead of the
`huaweicloud_networking_floatingip_associate_v2` resource.

## Example Usage

```hcl
resource "huaweicloud_networking_port" "myport" {
  network_id = "a5bbd213-e1d3-49b6-aed1-9df60ea94b9a"
}

resource "huaweicloud_vpc_eip" "myeip" {
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

resource "huaweicloud_networking_eip_associate" "associated" {
  public_ip = huaweicloud_vpc_eip.myeip.address
  port_id   = huaweicloud_networking_port.myport.id
}
```

## Argument Reference

The following arguments are supported:

* `public_ip` - (Required, String, ForceNew) The EIP to associate.

* `port_id` - (Required, String, ForceNew) ID of an existing port with at least one IP address to
    associate with this EIP.

## Attributes Reference

The following attributes are exported:


## Import

EIP associations can be imported using the `id` of the EIP, e.g.

```
$ terraform import huaweicloud_networking_eip_associate.eip 2c7f39f3-702b-48d1-940c-b50384177ee1
```
