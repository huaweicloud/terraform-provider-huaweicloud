---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_vpc_endpoint"
description: |-
  Manages the VPC endpoint via FunctionGraph side within HuaweiCloud.
---

# huaweicloud_fgs_vpc_endpoint

Manages the VPC endpoint via FunctionGraph side within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_fgs_vpc_endpoint" "test" {
  xrole     = var.iam_agency_name
  vpc_id    = var.vpc_id
  subnet_id = var.subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region where the VPC endpoint is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, NonUpdatable) The ID of the VPC to which the VPC endpoint belongs.

* `subnet_id` - (Required, String, NonUpdatable) The ID of the subnet to which the VPC endpoint belongs.

  -> This `subnet_id` must belong to the `vpc_id`.

* `flavor` - (Optional, String, NonUpdatable) The flavor of the VPC endpoint.  
  By default, the large specification will be used.

* `xrole` - (Optional, String, NonUpdatable) The IAM agency name of the VPC endpoint.  
  
  -> The agency must include these roles:<br>
     `vpc:zone:get`<br>
     `vpc:subnets:get`<br>
     `vpc:vpcs:get`<br>
     `vpc:securityGroups:get`<br>
     `vpc:securityGroups:list`<br>
     `vpcep:endpoints:create`<br>
     `vpcep:endpoints:delete`<br>
     `vpcep:endpoints:get`<br>
     `dns:zone:create`<br>
     `dns:zone:delete`<br>
     `dns:zone:get`<br>
     `dns:recordset:create`<br>
     `dns:recordset:update`<br>
     `dns:recordset:list`

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `endpoints` - The list of IP addresses of the VPC endpoint.

* `address` - The domain address of the endpoint.
