---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_gateway"
description: ""
---

# huaweicloud_nat_gateway

Use this data source to get an available public NAT gateway within HuaweiCloud.

## Example Usage

```hcl
variable "gateway_name" {}

data "huaweicloud_nat_gateway" "test" {
  name = var.gateway_name
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to create the NAT gateway resource. If omitted, the
  provider-level region will be used.

* `id` - (Optional, String) Specifies the ID of the NAT gateway.

* `name` - (Optional, String) Specifies the public NAT gateway name.  
  The valid length is limited from `1` to `64`, only letters, digits, hyphens (-) and underscores (_) are allowed.

* `subnet_id` - (Optional, String) Specifies the subnet ID of the downstream interface (the next hop of the DVR) of the
  public NAT gateway.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC this public NAT gateway belongs to.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the NAT gateway.

* `spec` - (Optional, String) The public NAT gateway type. The valid values are as follows:
  + **1**: Small type, which supports up to `10,000` SNAT connections.
  + **2**: Medium type, which supports up to `50,000` SNAT connections.
  + **3**: Large type, which supports up to `200,000` SNAT connections.
  + **4**: Extra-large type, which supports up to `1,000,000` SNAT connections.

* `description` - (Optional, String) Specifies the description of the NAT gateway. The value contains 0 to 255
  characters, and angle brackets (<)
  and (>) are not allowed.

* `status` - (Optional, String) Specifies the status of the NAT gateway.
