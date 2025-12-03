---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_loadbalancer_ports"
description: |-
  Use this data source to get the list of ports and IP addresses on the downstream subnet used by a load balancer.
---

# huaweicloud_elb_loadbalancer_ports

Use this data source to get the list of ports and IP addresses on the downstream subnet used by a load balancer.

## Example Usage

```hcl
variable "loadbalancer_id" {}

data "huaweicloud_elb_loadbalancer_ports" "test" {
  loadbalancer_id = var.loadbalancer_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `loadbalancer_id` - (Required, String) Specifies the load balancer ID.

* `port_id` - (Optional, List) Specifies the ID of the port used by the load balancer. Multiple IDs can be used.

* `ip_address` - (Optional, List) Specifies the private IPv4 address used by the load balancer. Multiple IP addresses
  can be used.

* `ipv6_address` - (Optional, List) Specifies the IPv6 address used by the load balancer. Multiple IPv6 addresses can be
  used.

* `type` - (Optional, List) Specifies the port type. Multiple types can be used. Value options:
  + **l4_hc**: port used for health checks during Layer 4 traffic forwarding using DNAT.
  + **l4**: port used for Layer 4 health checks and traffic forwarding using FullNAT.
  + **l7**: port used for Layer 7 health checks and service forwarding.

* `virsubnet_id` - (Optional, List) Specifies the network ID of the downstream subnet where the port is located. Multiple
  IDs can be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ports` - Indicates the list of ports used by the load balancer.
  The [ports](#ports_struct) structure is documented below.

<a name="ports_struct"></a>
The `ports` block supports:

* `port_id` - Indicates the ID of the port used by the load balancer.

* `ip_address` - Indicates the private IPv4 address used by the load balancer.

* `ipv6_address` - Indicates the IPv6 address used by the load balancer.

* `type` - Indicates the port type.

* `virsubnet_id` - Indicates the network ID of the downstream subnet where the port is located.
