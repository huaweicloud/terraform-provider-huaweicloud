---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_interface_attach"
description: ""
---

# huaweicloud_compute_interface_attach

Attaches a Network Interface to an Instance.

## Example Usage

### Attach a port (under the specified network) to the ECS instance and generate a random IP address

```hcl
variable "instance_id" {}
variable "network_id" {}

resource "huaweicloud_compute_interface_attach" "test" {
  instance_id = var.instance_id
  network_id  = var.network_id
}
```

### Attach a port (under the specified network) to the ECS instance and use the custom security groups

```hcl
variable "instance_id" {}
variable "network_id" {}
variable "security_group_ids" {
  type = list(string)
}

resource "huaweicloud_compute_interface_attach" "test" {
  instance_id        = var.instance_id
  network_id         = var.network_id
  security_group_ids = var.security_group_ids
}
```

### Attach a custom port to the ECS instance

```hcl
variable "security_group_id" {}

data "huaweicloud_vpc_subnet" "mynet" {
  name = "subnet-default"
}

data "huaweicloud_networking_port" "myport" {
  network_id = data.huaweicloud_vpc_subnet.mynet.id
  fixed_ip   = "192.168.0.100"
}

resource "huaweicloud_compute_instance" "myinstance" {
  name               = "instance"
  image_id           = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id          = "s6.small.1"
  key_pair           = "my_key_pair_name"
  security_group_ids = [var.security_group_id]
  availability_zone  = "cn-north-4a"

  network {
    uuid = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }
}

resource "huaweicloud_compute_interface_attach" "attached" {
  instance_id = huaweicloud_compute_instance.myinstance.id
  port_id     = data.huaweicloud_networking_port.myport.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the network interface attache resource. If
  omitted, the provider-level region will be used. Changing this creates a new network interface attache resource.

* `instance_id` - (Required, String, ForceNew) The ID of the Instance to attach the Port or Network to.

* `port_id` - (Optional, String, ForceNew) The ID of the Port to attach to an Instance.
  This option and `network_id` are mutually exclusive.

* `network_id` - (Optional, String) The ID of the Network to attach to an Instance. A port will be created automatically.
  This option and `port_id` are mutually exclusive.

* `fixed_ip` - (Optional, String) An IP address to associate with the port.

  ->This option cannot be used with `port_id`. You must specify a `network_id`. The IP address must lie in a range on
  the supplied network.

* `fixed_ipv6` - (Optional, String) The IPv6 address.

* `source_dest_check` - (Optional, Bool) Specifies whether the ECS processes only traffic that is destined specifically
  for it. This function is enabled by default but should be disabled if the ECS functions as a SNAT server or has a
  virtual IP address bound to it.

* `security_group_ids` - (Optional, List) Specifies the list of security group IDs bound to the specified port.  
  Defaults to the default security group.

* `delete_on_termination` - (Optional, String) Whether to delete a NIC when detaching it. Value options:
  + **true**: Delete the NIC.
  + **false**: Do not delete the NIC.

* `ipv6_enable` - (Optional, Bool, ForceNew) Specifies if the NIC supporting IPv6 or not.

* `ipv6_bandwidth_id` - (Optional, String, ForceNew) Specifies the shared bandwidth ID to which the IPv6 NIC attaches.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of ECS instance ID and port ID separated by a slash.

* `mac` - The MAC address of the NIC.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.

* `delete` - Default is 10 minutes.

## Import

Interface Attachments can be imported using the Instance ID and Port ID separated by a slash, e.g.

```shell
$ terraform import huaweicloud_compute_interface_attach.ai_1 89c60255-9bd6-460c-822a-e2b959ede9d2/45670584-225f-46c3-b33e-6707b589b666
```
