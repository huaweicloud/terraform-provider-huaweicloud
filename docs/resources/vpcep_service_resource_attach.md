---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_service_resource_attach"
description: -|
  Use this resource to manage a VPCEP endpoint service attach resource within HuaweiCloud.
---

# huaweicloud_vpcep_service_resource_attach

Use this resource to manage a VPCEP endpoint service attach resource within HuaweiCloud.

-> This resource is only a one-time action resource for adding server resources to VPCEP service.
   Deleting this resource will not remove the server resources from the VPCEP service, but will only remove the
   resource information from the tfstate file.

## Example Usage

```hcl
variable "service_id" {}
variable "server_resources" {
  type = list(object({
    resource_id          = string
    availability_zone_id = string
  }))
}

resource "huaweicloud_vpcep_service_resource_attach" "test" {
  service_id       = var.service_id
  server_resources = var.server_resources
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.  
  If omitted, the provider-level region will be used.  
  Change this parameter will create a new resource.

* `service_id` - (Required, String, NonUpdatable) Specifies the ID of the VPC endpoint service.

* `server_resources` - (Required, List, NonUpdatable) Specifies the list of server resources to be added to
  the VPC endpoint service.  
  The [server_resources](#vpcep_service_server_resources) structure is documented below.

<a name="vpcep_service_server_resources"></a>
The `server_resources` block supports:

* `resource_id` - (Required, String) Specifies the ID of the server resource.

* `availability_zone_id` - (Required, String) Specifies the availability zone ID of the server resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
