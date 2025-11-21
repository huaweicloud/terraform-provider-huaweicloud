---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_protectable_resources"
description: |-
  Use this data source to get the protectable resources of WAF within HuaweiCloud.
---

# huaweicloud_waf_protectable_resources

Use this data source to get the protectable resources of WAF within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_waf_protectable_resources" "test" {
  resource_type = "elb"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `resource_type` - (Required, String) Specifies the type of protection resource to query. Currently, only support **elb**.

* `vpc_id` - (Optional, String) Specifies the VPC ID where the load balancer is located.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The detailed information on protectable resources.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `loadbalancer_name` - The load balancer name.

* `loadbalancer_id` - The load balancer ID.

* `domain_id` - The account ID to which the load balancer belongs.

* `project_id` - The project ID.

* `listeners` - The current list of listeners associated with ELB.

  The [listeners](#items_listeners_struct) structure is documented below.

* `eips` - The EIP bound to the load balancer.

  The [eips](#items_eips_struct) structure is documented below.

<a name="items_listeners_struct"></a>
The `listeners` block supports:

* `name` - The name of the listener.

* `id` - The ID of the listener.

* `protocol` - The listening protocol of the listener.

* `protocol_port` - The listening port of the listener.

<a name="items_eips_struct"></a>
The `eips` block supports:

* `eip_id` - The elastic IP ID.

* `eip_address` - The elastic IP address.

* `ip_version` - The IP version.
