---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_private_dnat_rule"
description: |-
  Manages a DNAT rule resource of the **private** NAT within HuaweiCloud.
---

# huaweicloud_nat_private_dnat_rule

Manages a DNAT rule resource of the **private** NAT within HuaweiCloud.

## Example Usage

### DNAT rules forwarded with ECS instance as the backend

```hcl
variable "gateway_id" {}
variable "transit_ip_id" {}

resource "huaweicloud_compute_instance" "test" {
  ...
}

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id            = var.gateway_id
  protocol              = "tcp"
  transit_ip_id         = var.transit_ip_id
  transit_service_port  = 1000
  backend_interface_id  = huaweicloud_compute_instance.test.network[0].port
  internal_service_port = 2000
}
```

### DNAT rules forwarded with ELB loadbalancer as the backend

```hcl
variable "network_id" {}
variable "gateway_id" {}
variable "transit_ip_id" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  ...
}

data "huaweicloud_networking_port" "test" {
  network_id = var.network_id
  fixed_ip   = huaweicloud_elb_loadbalancer.test.ipv4_address
}

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id            = var.gateway_id
  protocol              = "tcp"
  transit_ip_id         = var.transit_ip_id
  transit_service_port  = 1000
  backend_interface_id  = data.huaweicloud_networking_port.test.id
  internal_service_port = 2000
}
```

### DNAT rules forwarded with VIP as the backend

```hcl
variable "network_id" {}
variable "gateway_id" {}
variable "transit_ip_id" {}

resource "huaweicloud_networking_vip" "test" {
  network_id = var.network_id
}

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id            = var.gateway_id
  protocol              = "tcp"
  transit_ip_id         = var.transit_ip_id
  transit_service_port  = 1000
  backend_interface_id  = huaweicloud_networking_vip.test.id
  internal_service_port = 2000
}
```

### DNAT rules forwarded with a custom private IP address as the backend

```hcl
variable "gateway_id" {}
variable "transit_ip_id" {}

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id            = var.gateway_id
  protocol              = "tcp"
  transit_ip_id         = var.transit_ip_id
  transit_service_port  = 1000
  backend_private_ip    = "172.168.0.69"
  internal_service_port = 2000
}
```

### DNAT rules for all ports

```hcl
variable "gateway_id" {}
variable "transit_ip_id" {}

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id         = var.gateway_id
  protocol           = "any"
  transit_ip_id      = var.transit_ip_id
  backend_private_ip = "172.168.0.69"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the DNAT rule is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `gateway_id` - (Required, String, ForceNew) Specifies the private NAT gateway ID to which the DNAT rule belongs.  
  Changing this will create a new resource.

* `transit_ip_id` - (Required, String) Specifies the ID of the transit IP for private NAT.

* `transit_service_port` - (Optional, Int) Specifies the port of the transit IP.  

-> Defaults to `0` and the default port is only available for rules with the protocol **any**.

* `protocol` - (Optional, String) Specifies the protocol type.  
  The valid values are **tcp**, **udp** and **any**. Defaults to **any**.

* `backend_interface_id` - (Optional, String) Specifies the network interface ID of the transit IP for private NAT.  
  Exactly one of `backend_interface_id` and `backend_private_ip` must be set.

* `backend_private_ip` - (Optional, String) Specifies the private IP address of the backend instance.

* `internal_service_port` - (Optional, Int) Specifies the port of the backend instance.

-> Defaults to `0` and the default port is only available for rules with the protocol **any**.

* `description` - (Optional, String) Specifies the description of the DNAT rule, which contain maximum of `255`
  characters, and angle brackets (< and >) are not allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `backend_type` - The type of backend instance.
  The valid values are as follows:
  + **COMPUTE**: ECS instance.
  + **VIP**: VIP.
  + **ELB**: ELB loadbalancer.
  + **ELBv3**: ver.3 ELB loadbalancer.
  + **CUSTOMIZE**: custom backend IP address.

* `created_at` - The creation time of the DNAT rule.

* `updated_at` - The latest update time of the DNAT rule.

* `enterprise_project_id` - The ID of the enterprise project to which the private DNAT rule belongs.

## Import

DNAT rules can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_nat_private_dnat_rule.test <id>
```
