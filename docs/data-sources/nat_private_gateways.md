---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_private_gateways"
description: ""
---

# huaweicloud_nat_private_gateways

Use this data source to get the list of private NAT gateways.

## Example Usage

```hcl
variable "gateway_name" {}

data "huaweicloud_nat_private_gateways" "test" {
  name = var.gateway_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the private NAT gateways are located.  
  If omitted, the provider-level region will be used.

* `gateway_id` - (Optional, String) Specifies the ID of the private NAT gateway.

* `name` - (Optional, String) Specifies the name of the private NAT gateway.  

* `spec` - (Optional, String) Specifies the specification of the private NAT gateways.  
  The valid values are as follows:
  + **Small**: Small type, which supports up to `20` rules, `200 Mbit/s` bandwidth, `20,000` PPS and `2,000` SNAT
    connections.
  + **Medium**: Medium type, which supports up to `50` rules, `500 Mbit/s` bandwidth, `50,000` PPS and `5,000` SNAT
    connections.
  + **Large**: Large type, which supports up to `200` rules, `2 Gbit/s` bandwidth, `200,000` PPS and `20,000` SNAT
    connections.
  + **Extra-Large**: Extra-large type, which supports up to `500` rules, `5 Gbit/s` bandwidth, `500,000` PPS and
    `50,000` SNAT connections.

* `status` - (Optional, String) Specifies the current status of the private NAT gateways.
  The valid values are as follows:
  + **ACTIVE**: The status of the private NAT gateway is normal operation.
  + **FROZEN**: The status of the private NAT gateway is frozen.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC to which the private NAT gateways belong.

* `subnet_id` - (Optional, String) Specifies the ID of the subnet to which the private NAT gateways belong.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the private NAT
  gateways belong.

* `description` - (Optional, List) Specifies the description of the private NAT gateway.  

* `tags` - (Optional, Map) Specifies the key/value pairs to associate the private NAT gateways.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `gateways` - The private NAT gateway list.
  The [gateways](#Gateway_Gateways) structure is documented below.

<a name="Gateway_Gateways"></a>
The `gateways` block supports:

* `id` - The ID of the private NAT gateway.

* `name` - The name of the private NAT gateway.  

* `spec` - The specification of the private NAT gateway.  
  The valid values are as follows:
  + **Small**: Small type, which supports up to `20` rules, `200 Mbit/s` bandwidth, `20,000` PPS and `2,000` SNAT
    connections.
  + **Medium**: Medium type, which supports up to `50` rules, `500 Mbit/s` bandwidth, `50,000` PPS and `5,000` SNAT
    connections.
  + **Large**: Large type, which supports up to `200` rules, `2 Gbit/s` bandwidth, `200,000` PPS and `20,000` SNAT
    connections.
  + **Extra-Large**: Extra-large type, which supports up to `500` rules, `5 Gbit/s` bandwidth, `500,000` PPS and
    `50,000` SNAT connections.

* `description` - The description of the private NAT gateway.

* `status` - The current status of the private NAT gateway.
  The valid values are as follows:
  + **ACTIVE**: The status of the private NAT gateway is normal operation.
  + **FROZEN**: The status of the private NAT gateway is frozen.

* `created_at` - The creation time of the private NAT gateway.

* `updated_at` - The latest update time of the private NAT gateway.

* `vpc_id` - The ID of the VPC to which the private NAT gateway belongs.

* `subnet_id` - The ID of the subnet to which the private NAT gateway belongs.

* `ngport_ip_address` - The IP address of the NG port of the private NAT gateway.

* `rule_max` - The maximum number of rules of the private NAT gateway.

* `enterprise_project_id` - The ID of the enterprise project to which the private NAT gateway belongs.

* `tags` - The key/value pairs to associate with the private NAT gateway.
