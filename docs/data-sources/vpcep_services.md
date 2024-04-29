---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_services"
description: ""
---

# huaweicloud_vpcep_services

Use this data source to get VPC endpoint services.

## Example Usage

```hcl
variable service_name {}
variable server_type {}
variable status {}

data "huaweicloud_vpcep_services" "services" {
  service_name = var.service_name
  server_type  = var.server_type
  status       = var.status
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the VPC endpoint services. If omitted, the
  provider-level region will be used.

* `service_id` - (Optional, String) Specifies the ID of VPC endpoint service.

* `service_name` - (Optional, String) Specifies the full name of the VPC endpoint service.

* `status` - (Optional, String) Specifies the status of the VPC endpoint service.
  The value can be **available** or **failed**.

* `server_type` - (Optional, String) Specifies the backend resource type. The valid values are as follows:
  + **VM**: Indicates the cloud server, which can be used as a server.
  + **LB**: Indicates the shared load balancer, which is applicable to services with high access traffic and services.

* `public_border_group` - (Optional, String) Specifies the information about Public Border Group of the pool
  corresponding to the VPC endpoint service.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `endpoint_services` - The list of VPC endpoint service.

The `endpoint_services` block supports:

* `id` - The ID of VPC endpoint service.

* `service_name` - The full name of the VPC endpoint service.

* `server_type` - The backend resource type.

* `vpc_id` - The ID of the VPC to which the backend resource of the VPC endpoint service belongs.

* `approval_enabled` - Whether connection approval is required. The default value is false.

* `service_type` - The type of VPC endpoint service.

* `created_at` - The creation time of VPC endpoint service.

* `updated_at` - The latest update time of VPC endpoint service.

* `connection_count` - The number of VPC endpoints that are in the Creating or Accepted status.

* `tcp_proxy` - Whether the client information, such as IP address, port number and marker_id, is
  transmitted to the server.

* `description` - The description of the VPC endpoint service.

* `enable_policy` - Whether the VPC endpoint policy is enabled. Defaults to **false**.

* `public_border_group` - The information about Public Border Group of the pool corresponding to
  the VPC endpoint service.

* `tags` - The key/value pairs to associate with the VPC endpoint service.
