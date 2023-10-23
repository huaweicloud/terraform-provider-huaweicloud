---
subcategory: Dedicated Load Balance (Dedicated ELB)
---

# huaweicloud_elb_loadbalancers

Use this data source to get the list of ELB loadbalancers.

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

* `name` - (Optional, String) Specifies the name of the ELB loadbalancer.

* `loadbalancer_id` - (Optional, String) Specifies the ID of the ELB loadbalancer.

* `description` - (Optional, String) Specifies the description of the ELB loadbalancer.

* `share_type` - (Optional, String) Specifies whether the load balancer is a dedicated load balancer, Value options:
  **dedicated**, **share**.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC where the load balancer resides.

* `ipv4_subnet_id` - (Optional, String) Specifies the ID of the IPv4 subnet where the load balancer resides.

* `ipv6_network_id` - (Optional, String) Specifies the ID of the port bound to the IPv6 address of the load balancer.

* `l4_flavor_id` - (Optional, String) Specifies the ID of a flavor at Layer 4.

* `l7_flavor_id` - (Optional, String) Specifies the ID of a flavor at Layer 7.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `loadbalancers` - Lists the load balancers.
  The [object](#loadbalancers_object) structure is documented below.

<a name="loadbalancers_object"></a>
The `loadbalancers` block supports:

* `id` - The loadbalancer ID.

* `name` - The loadbalancer name.

* `description` - The description of loadbalancer.

* `availability_zone` - The list of AZs where the load balancer is created..

* `cross_vpc_backend` - Whether to enable IP as a Backend Server.

* `vpc_id` - The ID of the VPC where the load balancer resides.

* `ipv4_subnet_id` - The  ID of the IPv4 subnet where the load balancer resides.

* `ipv6_network_id` - The ID of the IPv6 subnet where the load balancer resides.

* `ipv4_address` - The private IPv4 address bound to the load balancer.

* `ipv4_port_id` - The ID of the port bound to the private IPv4 address of the load balancer.

* `ipv6_address` - The IPv6 address bound to the load balancer.

* `l4_flavor_id` - The ID of a flavor at Layer 4.

* `l7_flavor_id` - The ID of a flavor at Layer 7

* `enterprise_project_id` - The enterprise project ID.

* `autoscaling_enabled` - Whether the current load balancer enables elastic expansion.

* `backend_subnets` - Lists the IDs of subnets on the downstream plane.

* `protection_status` - Modify the protection status, value: nonProtection: No protection, the default value is  
  nonProtection, consoleProtection: console modification protection.

* `protection_reason` - Reasons for setting protection, illustrate: Only valid when protection_status is
  consoleProtection, Minimum length=0, Maximum length=255.
