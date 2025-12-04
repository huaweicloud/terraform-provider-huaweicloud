---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_loadbalancer_status"
description: |-
  Use this data source to get the status tree of a load balancer and to show information about all resources associated
  with the load balancer.
---

# huaweicloud_elb_loadbalancer_status

Use this data source to get the status tree of a load balancer and to show information about all resources associated
with the load balancer.

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

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `statuses` - Indicates the supplementary information about the load balancer status tree.
  The [statuses](#statuses_struct) structure is documented below.

<a name="statuses_struct"></a>
The `statuses` block supports:

* `loadbalancer` - Indicates the statuses of the load balancer and its associated resources.
  The [loadbalancer](#loadbalancer_struct) structure is documented below.

<a name="loadbalancer_struct"></a>
The `loadbalancer` block supports:

* `id` - Indicates the load balancer ID.

* `name` - Indicates the load balancer name.

* `provisioning_status` - Indicates the provisioning status of the load balancer.

* `listeners` - Indicates the listeners added to the load balancer.
  The [listeners](#listeners_struct) structure is documented below.

* `pools` - Indicates the backend server groups associated with the load balancer.
  The [pools](#pools_struct) structure is documented below.

* `provisioning_status` - Indicates the operating status of the load balancer.

<a name="listeners_struct"></a>
The `listeners` block supports:

* `id` - Indicates the listener ID.

* `name` - Indicates the name of the listener.

* `provisioning_status` - Indicates the provisioning status of the listener.

* `pools` - Indicates the operating status of the backend server groups associated with the listener.
  The [pools](#pools_struct) structure is documented below.

* `l7policies` - Indicates the operating status of the forwarding policy added to the listener.
  The [l7policies](#l7policies_struct) structure is documented below.

* `operating_status` - Indicates the operating status of the listener.

<a name="pools_struct"></a>
The `pools` block supports:

* `id` - Indicates the backend server group ID.

* `name` - Indicates the backend server group name.

* `provisioning_status` - Indicates the provisioning status of the backend server group.

* `healthmonitor` - Indicates the health check results of backend servers in the load balancer status tree.
  The [healthmonitor](#healthmonitor_struct) structure is documented below.

* `members` - Indicates the statuses of all backend servers in the backend server group.
  The [members](#members_struct) structure is documented below.

* `operating_status` - Indicates the operating status of the backend server group.

<a name="healthmonitor_struct"></a>
The `healthmonitor` block supports:

* `id` - Indicates the health check ID.

* `name` - Indicates the health check name.

* `type` - Indicates the health check protocol.

* `provisioning_status` - Indicates the provisioning status of the health check.

<a name="members_struct"></a>
The `members` block supports:

* `id` - Indicates the backend server ID.

* `address` - Indicates the IP address of the backend server.

* `provisioning_status` - Indicates the provisioning status of the backend server.

* `protocol_port` - Indicates the port used by the backend server to receive requests.

* `operating_status` - Indicates the operating status of the backend server.

<a name="l7policies_struct"></a>
The `l7policies` block supports:

* `id` - Indicates the forwarding policy ID.

* `name` - Indicates the forwarding policy name.

* `action` - Indicates whether requests are forwarded to another backend server group or redirected to an HTTPS listener.

* `provisioning_status` - Indicates the provisioning status of the forwarding policy.

* `rules` - Indicates the status of all forwarding rules in the forwarding policy.
  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - Indicates the ID of the forwarding rule.

* `type` - Indicates the type of the match content.

* `provisioning_status` - Indicates the provisioning status of the forwarding rule.
