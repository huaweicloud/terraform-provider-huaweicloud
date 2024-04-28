---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_gateways"
description: ""
---

# huaweicloud_nat_gateways

Use this data source to get the list of NAT gateways.

## Example Usage

```hcl
variable "gateway_name" {}

data "huaweicloud_nat_gateways" "test" {
  name = var.gateway_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the NAT gateways are located.  
  If omitted, the provider-level region will be used.

* `gateway_id` - (Optional, String) Specifies the ID of the NAT gateway.

* `name` - (Optional, String) Specifies the name of the NAT gateway.

* `spec` - (Optional, String) Specifies the specification of the NAT gateways.
  The valid values are as follows:
  + **1**: Small type, which supports up to `10,000` SNAT connections.
  + **2**: Medium type, which supports up to `50,000` SNAT connections.
  + **3**: Large type, which supports up to `200,000` SNAT connections.
  + **4**: Extra-large type, which supports up to `1,000,000` SNAT connections.

* `status` - (Optional, String) Specifies the current status of the NAT gateways.
  The valid values are as follows:
  + **ACTIVE**: The status of the NAT gateway is available.
  + **INACTIVE**: The status of the NAT gateway is unavailable.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC to which the NAT gateways belong.

* `subnet_id` - (Optional, String) Specifies the network ID of the downstream interface (the next hop of the DVR) of
  the NAT gateways.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the NAT
  gateways belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `gateways` - The list of the NAT gateway.
  The [gateways](#GatewayPublicGateways) structure is documented below.

<a name="GatewayPublicGateways"></a>
The `gateways` block supports:

* `id` - The ID of the NAT gateway.

* `name` - The name of the NAT gateway.

* `spec` - The specification of the NAT gateway.
  The valid values are as follows:
  + **1**: Small type, which supports up to `10,000` SNAT connections.
  + **2**: Medium type, which supports up to `50,000` SNAT connections.
  + **3**: Large type, which supports up to `200,000` SNAT connections.
  + **4**: Extra-large type, which supports up to `1,000,000` SNAT connections.

* `status` - The current status of the NAT gateway.
  The valid values are as follows:
  + **ACTIVE**: The status of the NAT gateway is available.
  + **INACTIVE**: The status of the NAT gateway is unavailable.

* `description` - The description of the NAT gateway.

* `created_at` - The creation time of the NAT gateway.

* `vpc_id` - The ID of the VPC to which the NAT gateway belongs.

* `subnet_id` - The network ID of the downstream interface (the next hop of the DVR) of the NAT gateway.

* `enterprise_project_id` - The ID of the enterprise project to which the NAT gateway belongs.
