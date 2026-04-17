---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_service_add_servers"
description: -|
  Use this resource to add backend resources to a VPC endpoint service within HuaweiCloud.
---

# huaweicloud_vpcep_service_add_servers

Use this resource to add backend resources to a VPC endpoint service within HuaweiCloud.

## Example Usage

```hcl
variable "vpc_endpoint_service_id" {}
variable "elb_instance_id_1" {}
variable "elb_az_id_1" {}
variable "elb_instance_id_2" {}
variable "elb_az_id_2" {}

resource "huaweicloud_vpcep_service_add_servers" "test" {
  vpc_endpoint_service_id = var.vpc_endpoint_service_id
  
  server_resources {
    resource_id          = var.elb_instance_id_1
    availability_zone_id = var.elb_az_id_1
  }

  server_resources {
    resource_id          = var.elb_instance_id_2
    availability_zone_id = var.elb_az_id_2
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the VPC endpoint service.
  If omitted, the provider-level region will be used. Changing this creates a new VPC endpoint service resource.

* `vpc_endpoint_service_id` - (Required, String, NonUpdatable) Specifies the ID of the VPC endpoint service.

* `server_resources` - (Required, List, NonUpdatable) Specifies the load balancer IDs and AZs.  
  The [server_resources](#server_resources) structure is documented below.

<a name="server_resources"></a>
The `server_resources` block supports:

* `resource_id` - (Required, String, NonUpdatable) Specifies the load balancer ID.

* `availability_zone_id` - (Optional, String, NonUpdatable) Specifies the ID of the AZ where the load balancer is located.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the VPC endpoint service.
