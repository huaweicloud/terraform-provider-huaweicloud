---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_endpoints"
description: ""
---

# huaweicloud_vpcep_endpoints

Use this data source to get a list of VPC endpoints.

## Example Usage

```hcl
variable "service_name" {}
variable "endpoint_id" {}

data "huaweicloud_vpcep_endpoints" "endpoints" {
  service_name = var.service_name
  endpoint_id  = var.endpoint_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the VPC endpoints.
  If omitted, the provider-level region will be used.

* `service_name` - (Optional, String) The name of the VPC endpoint service.

* `vpc_id` - (Optional, String) The ID of the VPC where the endpoint is created.

* `endpoint_id` - (Optional, String) The ID of the VPC endpoint.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `endpoints` - A list of VPC endpoints.

The `endpoints` block supports:

* `id` - The ID of the VPC endpoint.

* `service_type` - The type of the VPC endpoint service. The valid value can be:
  + **gateway**: Configured by operation and maintenance personnel.
  Users do not need to create it and can use it directly.
  + **interface**: Including cloud services configured by operation and maintenance personnel
  and private services created by users themselves. Among them, the cloud services configured
  by operation and maintenance personnel do not need to be created, and users can use them directly.
  You can query the public endpoint service list to view all user-visible and connectable
  endpoint services configured by operation and maintenance personnel, and create an **interface**
  type endpoint service by creating an endpoint service.

* `status` - The status of the VPC endpoint. The value can be **accepted**, **pendingAcceptance**, **creating**,
  **failed**, **deleting**, or **rejected**.

* `service_name` - The name of the VPC endpoint service.

* `service_id` - The ID of the VPC endpoint service.

* `packet_id` - The packet ID of the VPC endpoint.
  OBS uses this field to implement double-ended fixed features, and distinguish what endpoint instance it is.

* `enable_dns` - Whether to create a private domain name. The value can be **true** or **false**.
  When the type of endpoint service is **gateway**, the field is not valid.

* `ip_address` - The IP address for accessing the associated VPC endpoint service.

* `description` - The description of VPC endpoint.

* `vpc_id` - The ID of the VPC where the endpoint is created.

* `subnet_id` - The network ID of the subnet in the VPC specified by `vpc_id`, in UUID format.

* `created_at` - The creation time of the VPC endpoint. Use UTC time format, the format is: YYYY-MM-DDTHH:MM:SSZ.

* `updated_at` - The update time of the VPC endpoint. Use UTC time format, the format is: YYYY-MM-DDTHH:MM:SSZ.

* `tags` - The key/value pairs associating with the VPC endpoint.

* `whitelist` - The list of IP address or CIDR block which can be accessed to the
  VPC endpoint.
  If not created, an empty list is returned.
  This field is displayed when creating an endpoint connected to an **interface** type endpoint service.
  Array length: **1** - **1000**.

* `enable_whitelist` - Whether to enable access control.
  This field is displayed when creating an endpoint connected to an **interface** type endpoint service.
