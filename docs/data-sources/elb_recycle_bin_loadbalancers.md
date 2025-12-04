---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_recycle_bin_loadbalancers"
description: |-
  Use this data source to get the list of load balancers in the recycle bin.
---

# huaweicloud_elb_recycle_bin_loadbalancers

Use this data source to get the list of load balancers in the recycle bin.

## Example Usage

```hcl
data "huaweicloud_elb_recycle_bin_loadbalancers" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `loadbalancer_id` - (Optional, List) Specifies the load balancer ID.
  Multiple IDs can be used.

* `name` - (Optional, List) Specifies the load balancer name.
  Multiple names can be used.

* `description` - (Optional, List) Specifies the load balancer description.
  Multiple descriptions can be used.

* `operating_status` - (Optional, List) Specifies the operating status of the load balancer.
  Multiple operating statuses can be used. Value options:
  + **ONLINE**: The load balancer is working normally.
  + **FROZEN**: The load balancer has been frozen.

* `guaranteed` - (Optional, String) Specifies whether the load balancer is a dedicated load balancer. Value options:
  **false**: The load balancer is a shared load balancer.
  + **true**: The load balancer is a dedicated load balancer.

* `vpc_id` - (Optional, List) Specifies the ID of the VPC where the load balancer resides.
  Multiple IDs can be used.

* `vip_port_id` - (Optional, List) Specifies the ID of the port bound to the private IPv4 address of the load balancer.
  Multiple IDs can be used.

* `vip_address` - (Optional, List) Specifies the IPv4 private IP address of the load balancer.
  Multiple IP addresses can be used.

* `vip_subnet_cidr_id` - (Optional, List) Specifies the ID of the IPv4 subnet where the load balancer resides.
  Multiple IDs can be used.

* `ipv6_vip_port_id` - (Optional, List) Specifies the ID of the port bound to the IPv6 address of the load balancer.
  Multiple ports can be used.

* `ipv6_vip_address` - (Optional, List) Specifies the IPv6 address bound to the load balancer.
  Multiple IPv6 addresses can be used.

* `ipv6_vip_virsubnet_id` - (Optional, List) Specifies the ID of the IPv6 subnet where the load balancer resides.
  Multiple IDs can be used.

* `publicips` - (Optional, List) Specifies the EIP bound to the load balancer.
  Multiple publicips can be used.
  + If `publicip_id` is used as the query condition, the format is **publicips=publicip_id=xxx**.
  + If `publicip_address` is used as the query condition, the format is **publicips=publicip_address=xxx**.
  + If `ip_version` is used as the query condition, the format is **publicips=ip_version=xxx**.

* `availability_zone_list` - (Optional, List) Specifies the list of AZs where the load balancer is created.
  Multiple AZs can be used.

* `l4_flavor_id` - (Optional, List) Specifies the ID of a flavor at Layer 4.
  Multiple IDs can be used.

* `l7_flavor_id` - (Optional, List) Specifies the ID of a flavor at Layer 7.
  Multiple flavors can be used.

* `billing_info` - (Optional, List) Specifies the provides resource billing information.
  Multipple values can be used.

* `member_device_id` - (Optional, List) Specifies the ID of the ECS that is associated with the load balancer as a
  backend server. Multiple IDs can be used.

* `member_address` - (Optional, List) Specifies the private IP address of the ECS that is associated with the load
  balancer as a backend server. Multiple private IP addresses can be used.

* `enterprise_project_id` - (Optional, List) Specifies the ID of the enterprise project where the load balancer is used.
  Multiple values can be used.

* `ip_version` - (Optional, List) Specifies the IP version.
  Multiple versions can be used. The value can be **4 (IPv4)** or **6 (IPv6)**.

* `deletion_protection_enable` - (Optional, String) Specifies whether to enable deletion protection. Value options:
  + **false**: disable this option
  + **true**: enable this option

* `elb_virsubnet_type` - (Optional, List) Specifies the type of the subnet on the downstream plane.
  Multiple values can be used. Value options:
  + **ipv4**: IPv4 subnet
  + **dualstack**: subnet that supports IPv4/IPv6 dual stack

* `protection_status` - (Optional, List) Specifies the protection status. Value options:
  + **nonProtection**: The resource is not protected.
  + **consoleProtection**: Modification is not allowed on the console.

* `global_eips` - (Optional, List) Specifies the EIP bound to the load balancer.
  Multiple values can be used. Value options:
  + If `global_eip_id` is used as the query condition, the format is **global_eips=global_eip_id=xxx**.
  + If `global_eip_address` is used as the query condition, the format is **global_eips=global_eip_address=xxx**.
  + If `ip_version` is used as the query condition, the format is **global_eips=ip_version=xxx**.

* `log_topic_id` - (Optional, String) Specifies the ID of the log topic that is associated with the load balancer.
  Multiple IDs can be used.

* `log_group_id` - (Optional, String) Specifies the ID of the log group that is associated with the load balancer.
  Multiple IDs can be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `loadbalancers` - Indicates the list of the load balancers in the recycle bin.

  The [loadbalancers](#loadbalancers_struct) structure is documented below.

<a name="loadbalancers_struct"></a>
The `loadbalancers` block supports:

* `id` - Indicates the load balancer ID.

* `name` - Indicates the load balancer name.

* `availability_zone_list` - Indicates the list of AZs where the load balancers are created.

* `description` - Indicates the provides supplementary information about the load balancer.

* `vpc_id` - Indicates the ID of the VPC where the load balancer is located.

* `l4_flavor_id` - Indicates the ID of a flavor at Layer 4.

* `l4_scale_flavor_id` - Indicates the ID of the elastic flavor at Layer 4.

* `l7_flavor_id` - Indicates the ID of a flavor at Layer 7.

* `l7_scale_flavor_id` - Indicates the ID of an elastic flavor at Layer 7.

* `ipv6_vip_virsubnet_id` - Indicates the ID of the IPv6 subnet where the load balancer is located.

* `ipv6_vip_address` - Indicates the IPv6 address of the load balancer.

* `ip_target_enable` - Indicates whether to add backend servers that are not in the load balancer's VPC.

* `pools` - Indicates the IDs of backend server groups associated with the load balancer.

  The [pools](#loadbalancers_pools_struct) structure is documented below.

* `global_eips` - Indicates the global EIP bound to the load balancer.

  The [global_eips](#loadbalancers_global_eips_struct) structure is documented below.

* `frozen_scene` - Indicates the scenario where the load balancer is frozen.

* `ipv6_bandwidth` - Indicates the ID of the bandwidth.

  The [ipv6_bandwidth](#loadbalancers_ipv6_bandwidth_struct) structure is documented below.

* `provider` - Indicates the provider of the load balancer. The value is fixed to **vlb**.

* `protection_status` - Indicates the provisioning status of the load balancer.

* `log_group_id` - Indicates the ID of the log group that is associated with the load balancer.

* `vip_address` - Indicates the private IPv4 address of the load balancer.

* `vip_port_id` - Indicates the ID of the port bound to the private IPv4 address of the load balancer.

* `publicips` - Indicates the EIP bound to the load balancer.

  The [publicips](#loadbalancers_publicips_struct) structure is documented below.

* `charge_mode` - Indicates the charge mode of the load balancer.

* `operating_status` - Indicates the operating status of the load balancer.

* `enterprise_project_id` - Indicates the ID of the enterprise project.

* `deletion_protection_enable` - Indicates whether to enable deletion protection.

* `provisioning_status` - Indicates the provisioning status of the load balancer.

* `elb_virsubnet_ids` - Indicates the network IDs of subnets on the downstream plane.

* `public_border_group` - Indicates the AZ group to which the load balancer belongs.

* `waf_failure_action` - Indicates traffic distributing policies when the WAF is faulty.

* `ipv6_vip_port_id` - Indicates the list of AZs where the load balancers are created.

* `guaranteed` - Indicates whether the load balancer is a dedicated load balancer.

* `billing_info` - Indicates the povides resource billing information.

* `elb_virsubnet_type` - Indicates the type of the subnet on the downstream plane.

* `protection_reason` - Indicates why modification protection is enabled.

* `log_topic_id` - Indicates the ID of the log topic that is associated with the load balancer.

* `listeners` - Indicates the IDs of listeners associated with the load balancer.

  The [listeners](#loadbalancers_listeners_struct) structure is documented below.

* `vip_subnet_cidr_id` - Indicates the ID of the IPv4 subnet where the load balancer is located.

* `tags` - Indicates the tags added to the load balancer.

  The [tags](#loadbalancers_tags_struct) structure is documented below.

* `auto_terminate_time` - Indicates the time when the load balancers in the recycle bin will be permanently deleted.
  The format is **yyyy-MM-dd'T'HH:mm:ss'Z'**.

* `created_at` - Indicates the time when the load balancer was created.

* `updated_at` - Indicates the time when the load balancer was updated.

<a name="loadbalancers_pools_struct"></a>
The `pools` block supports:

* `id` - Indicates the backend server group ID.

<a name="loadbalancers_global_eips_struct"></a>
The `global_eips` block supports:

* `global_eip_id` - Indicates the ID of the global EIP.

* `global_eip_address` - Indicates the global EIP.

* `ip_version` - Indicates the IP address version.

<a name="loadbalancers_ipv6_bandwidth_struct"></a>
The `ipv6_bandwidth` block supports:

* `id` - Indicates the ID of the shared bandwidth.

<a name="loadbalancers_publicips_struct"></a>
The `publicips` block supports:

* `publicip_id` - Indicates the EIP ID.

* `publicip_address` - Indicates the EIP.

* `ip_version` - Indicates the IP address version.

<a name="loadbalancers_listeners_struct"></a>
The `listeners` block supports:

* `id` - Indicates the listener ID.

<a name="loadbalancers_tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the tag key.

* `value` - Indicates the tag value.
