---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_loadbalancers"
description: |-
  Use this data source to get the list of ELB load balancers.
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

* `availability_zone` - (Optional, String) Specifies the list of AZ where the load balancer is created.

* `billing_info` - (Optional, String) Specifies the provides resource billing information.

* `deletion_protection_enable` - (Optional, String) Specifies whether the deletion protection is enabled. Value options:
  **true**, **false**.

* `global_eips` - (Optional, List) Specifies the global EIPs bound to the load balancer. It can be queried by different
  conditions:
  + If `global_eip_id` is used as the query condition, the format is **global_eip_id=xxx**
  + If `global_eip_address` is used as the query condition, the format is **global_eip_address=xxx**
  + If `ip_version` is used as the query condition, the format is **ip_version=xxx**

* `ipv6_address` - (Optional, String) Specifies the IPv6 address bound to the load balancer.

* `ipv6_vip_port_id` - (Optional, String) Specifies the ID of the port bound to the IPv6 address of the load balancer.

* `log_group_id` - (Optional, String) Specifies the ID of the log group that is associated with the load balancer.

* `log_topic_id` - (Optional, String) Specifies the ID of the log topic that is associated with the load balancer.

* `member_address` - (Optional, String) Specifies the private IP address of the cloud server that is associated with the
  load balancer as a backend server.

* `member_device_id` - (Optional, String) Specifies the ID of the cloud server that is associated with the load balancer
  as a backend server.

* `operating_status` - (Optional, String) Specifies the operating status of the load balancer. Value options:
  + **ONLINE**: indicates that the load balancer is running normally.
  + **FROZEN**: indicates that the load balancer is frozen.

* `protection_status` - (Optional, String) Specifies the protection status. Value options:
  + **nonProtection**: The load balancer is not protected.
  + **consoleProtection**: Modification Protection is enabled on the console.

* `provisioning_status` - (Optional, String) Specifies the provisioning status of the load balancer. Value options:
  + **ACTIVE**: The load balancer is successfully provisioned.
  + **PENDING_DELETE**: The load balancer is being deleted.

* `publicips` - (Optional, List) Specifies the EIPs bound to the load balancer. It can be queried by different conditions:
  + If `publicip_id` is used as the query condition, the format is **publicip_id=xxx**
  + If `publicip_address` is used as the query condition, the format is **publicip_address=xxx**
  + If `ip_version` is used as the query condition, the format is **ip_version=xxx**

* `ipv4_address` - (Optional, String) Specifies the private IPv4 address bound to the load balancer.

* `ipv4_port_id` - (Optional, String) Specifies the ID of the port bound to the private IPv4 address of the load balancer.

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

* `billing_info` - The provides resource billing information.

* `charge_mode` - The charge mode when of the load balancer.

* `deletion_protection_enable` - Whether the deletion protection is enabled.

* `frozen_scene` - The scenario where the load balancer is frozen.

* `global_eips` - The EIPs bound to the load balancer.
  The [global_eips](#global_eips_struct) structure is documented below.

* `listeners` - The list of listeners added to the load balancer.
  The [listeners](#listeners_struct) structure is documented below.

* `pools` - The list of pools associated with the load balancer.
  The [pools](#pools_struct) structure is documented below.

* `log_group_id` - The ID of the log group that is associated with the load balancer.

* `log_topic_id` - The ID of the log topic that is associated with the load balancer.

* `operating_status` - The operating status of the load balancer.

* `provisioning_status` - The provisioning status of the load balancer.

* `public_border_group` - The AZ group to which the load balancer belongs.

* `publicips` - The EIPs bound to the load balancer.
  The [publicips](#publicips_struct) structure is documented below.

* `waf_failure_action` - The traffic distributing policies when the WAF is faulty.

* `cross_vpc_backend` - Whether to enable IP as a Backend Server.

* `vpc_id` - The ID of the VPC where the load balancer resides.

* `ipv4_subnet_id` - The  ID of the IPv4 subnet where the load balancer resides.

* `ipv6_network_id` - The ID of the IPv6 subnet where the load balancer resides.

* `ipv4_address` - The private IPv4 address bound to the load balancer.

* `ipv4_port_id` - The ID of the port bound to the private IPv4 address of the load balancer.

* `ipv6_address` - The IPv6 address bound to the load balancer.

* `ipv6_vip_port_id` - The ID of the port bound to the IPv6 address of the load balancer.

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

* `created_at` - The time when the load balancer was created.

* `updated_at` - The time when the load balancer was updated.

<a name="global_eips_struct"></a>
The `global_eips` block supports:

* `global_eip_address` - The global EIP address

* `global_eip_id` - The ID of the global EIP.

* `ip_version` - The IP version.

<a name="listeners_struct"></a>
The `listeners` block supports:

* `id` - The listener ID.

<a name="pools_struct"></a>
The `pools` block supports:

* `id` - The pool ID.

<a name="publicips_struct"></a>
The `publicips` block supports:

* `publicip_id` - The EIP ID.

* `publicip_address` - The IP address.

* `ip_version` - The IP version.
