---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_private_gateway"
description: |-
  Manages a gateway resource of the **private** NAT within HuaweiCloud.
---

# huaweicloud_nat_private_gateway

Manages a gateway resource of the **private** NAT within HuaweiCloud.

## Example Usage

```hcl
variable "subnet_id" {}
variable "gateway_name" {}
variable "gateway_spec" {}

resource "huaweicloud_nat_private_gateway" "test" {
  subnet_id             = var.subnet_id
  name                  = var.gateway_name
  spec                  = var.gateway_spec
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the private NAT gateway is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the network ID of the subnet to which the private NAT gateway
  belongs.  
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the private NAT gateway name.  
  The valid length is limited from `1` to `64`, only English letters, Chinese characters, digits, hyphens (-) and
  underscores (_) are allowed.

* `spec` - (Optional, String) Specifies the specification of the private NAT gateway.  
  The valid values are as follows:
  + **Small**: Small type, which supports up to `20` rules, `200 Mbit/s` bandwidth, `20,000` PPS and `2,000` SNAT
    connections.
  + **Medium**: Medium type, which supports up to `50` rules, `500 Mbit/s` bandwidth, `50,000` PPS and `5,000` SNAT
    connections.
  + **Large**: Large type, which supports up to `200` rules, `2 Gbit/s` bandwidth, `200,000` PPS and `20,000` SNAT
    connections.
  + **Extra-Large**: Extra-large type, which supports up to `500` rules, `5 Gbit/s` bandwidth, `500,000` PPS and
    `50,000` SNAT connections.

  Defaults to **Small**.

* `description` - (Optional, String) Specifies the description of the private NAT gateway, which contain maximum of
  `255` characters, and angle brackets (< and >) are not allowed.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the ID of the enterprise project to which the private
  NAT gateway belongs.  
  Changing this will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the private NAT gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `created_at` - The creation time of the private NAT gateway.

* `updated_at` - The latest update time of the private NAT gateway.

* `status` - The current status of the private NAT gateway.

* `vpc_id` - The ID of the VPC to which the private NAT gateway belongs.

## Import

The private NAT gateways can be imported using their `id`, e.g.

```bash
$ terraform import huaweicloud_nat_private_gateway.test <id>
```
