---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud_networking_vip

Manages a VIP resource within HuaweiCloud. This is an alternative to `huaweicloud_networking_vip_v2`

## Example Usage

```hcl
data "huaweicloud_vpc_subnet" "mynet" {
  name = "subnet-default"
}

resource "huaweicloud_networking_vip" "myvip" {
  network_id = data.huaweicloud_vpc_subnet.mynet.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the vip resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `network_id` - (Required, String, ForceNew) Specifies the ID of the network to which the vip belongs.
  Changing this creates a new resource.

* `ip_version` - (Optional, Int, ForceNew) Specifies the IP version, either 4 (default) or 6.
  Changing this creates a new resource.

* `ip_address` - (Optional, String, ForceNew) Specifies the IP address desired in the subnet for this vip.
  Changing this creates a new resource.

* `name` - (Optional, String) Specifies a unique name for the vip.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the vip.
* `mac_address` - The MAC address of the vip.
* `status` - The status of vip.
* `device_owner` - The device owner of the vip.
* `subnet_id` - The subnet ID in which to allocate IP address for this vip.

## Import

Networking VIP can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_networking_vip.myvip ce595799-da26-4015-8db5-7733c6db292e
```
