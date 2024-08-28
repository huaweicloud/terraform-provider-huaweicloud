---
subcategory: Dedicated Load Balance (Dedicated ELB)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_loadbalancers"
description: ""
---

# huaweicloud_elb_loadbalancers

Use this data source to get the list of ELB load balancers.

## Example Usage

```hcl
variable "loadbalancer_name" {}

data "huaweicloud_elb_loadbalancers" "test" {
  name = var.loadbalancer_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the ELB load balancer.

* `loadbalancer_id` - (Optional, String) Specifies the ID of the ELB load balancer.

* `description` - (Optional, String) Specifies the description of the ELB load balancer.

* `type` - (Optional, String) Specifies whether the load balancer is a dedicated load balancer, Value options:
  **dedicated**, **share**.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC where the load balancer resides.

* `ipv4_subnet_id` - (Optional, String) Specifies the ID of the IPv4 subnet where the load balancer resides.

* `ipv6_network_id` - (Optional, String) Specifies the ID of the port bound to the IPv6 address of the load balancer.

* `l4_flavor_id` - (Optional, String) Specifies the ID of a flavor at Layer 4.

* `l7_flavor_id` - (Optional, String) Specifies the ID of a flavor at Layer 7.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `loadbalancers` - The List of load balancers.
  The [loadbalancers](#Elb_loadbalancer_loadbalancers) structure is documented below.

<a name="Elb_loadbalancer_loadbalancers"></a>
The `loadbalancers` block supports:

* `id` - The load balancer ID.

* `name` - The load balancer name.

* `loadbalancer_type` - The type of the load balancer.

* `description` - The description of load balancer.

* `availability_zone` - The list of AZs where the load balancer is created.

* `cross_vpc_backend` - Whether to enable IP as a Backend Server.

* `vpc_id` - The ID of the VPC where the load balancer resides.

* `ipv4_subnet_id` - The  ID of the IPv4 subnet where the load balancer resides.

* `ipv6_network_id` - The ID of the IPv6 subnet where the load balancer resides.

* `ipv4_address` - The private IPv4 address bound to the load balancer.

* `ipv4_port_id` - The ID of the port bound to the private IPv4 address of the load balancer.

* `ipv6_address` - The IPv6 address bound to the load balancer.

* `l4_flavor_id` - The ID of a flavor at Layer 4.

* `l7_flavor_id` - The ID of a flavor at Layer 7.

* `gw_flavor_id` - The flavor ID of the gateway load balancer.

* `min_l7_flavor_id` - The minimum seven-layer specification ID (specification type L7_elastic) for elastic expansion
  and contraction

* `enterprise_project_id` - The enterprise project ID.

* `autoscaling_enabled` - Whether the current load balancer enables elastic expansion.

* `backend_subnets` - Lists the IDs of subnets on the downstream plane.

* `protection_status` - The protection status for update.

* `protection_reason` - The reason for update protection.

* `type` - Whether the load balancer is a dedicated load balancer.
