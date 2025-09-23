---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_private_transit_ips"
description: |-
  Use this data source to get the list of transit IPs.
---

# huaweicloud_nat_private_transit_ips

Use this data source to get the list of transit IPs.

## Example Usage

```hcl
variable "ip_address" {}

data "huaweicloud_nat_private_transit_ips" "test" {
  ip_address = var.ip_address
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the transit IPs are located.  
  If omitted, the provider-level region will be used.

* `transit_ip_id` - (Optional, String) Specifies the ID of the transit IP.

* `ip_address` - (Optional, String) Specifies the IP address of the transit IP.  

* `gateway_id` - (Optional, String) Specifies the ID of the private NAT gateway to which the transit IP belongs.

* `subnet_id` - (Optional, String) Specifies the ID of the subnet to which the transit IPs belong.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate the transit IPs used for filter.

* `network_interface_id` - (Optional, String) Specifies the network interface ID of the transit IP for private NAT.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the transit
  IPs belong.

* `transit_subnet_id` - (Optional, List) Specifies the ID of the transit subnet.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `transit_ips` - The list ot the transit IPs.
  The [transit_ips](#private_transitIps) structure is documented below.

<a name="private_transitIps"></a>
The `transit_ips` block supports:

* `id` - The ID of the transit IP.

* `ip_address` - The IP address of the transit IP.

* `subnet_id` - The ID of the subnet to which the transit IP belongs.

* `tags` - The key/value pairs to associate with the transit IP.

* `gateway_id` - The ID of the private NAT gateway to which the transit IP belongs.

* `status` - The status of the transit IP.

* `created_at` - The creation time of the transit IP.

* `updated_at` - The latest update time of the transit IP.

* `network_interface_id` - The network interface ID of the transit IP for private NAT.

* `enterprise_project_id` - The ID of the enterprise project to which the transit IP belongs.
