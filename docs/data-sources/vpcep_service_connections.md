---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_service_connections"
description: ""
---

# huaweicloud_vpcep_service_connections

Use this data source to get VPC endpoint service connections.

## Example Usage

```hcl
variable service_id {}
variable endpoint_id {}
variable status {}

data "huaweicloud_vpcep_service_connections" "connections" {
  service_id  = var.service_id
  endpoint_id = var.endpoint_id
  status      = var.status
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the VPC endpoint service. If omitted, the
  provider-level region will be used.

* `service_id` - (Required, String) Specifies the ID of VPC endpoint service.

* `endpoint_id` - (Optional, String) Specifies the ID of VPC endpoint which has connected to
  VPC endpoint service.

* `status` - (Optional, String) Specifies the connection status of the VPC endpoint.
  The value can be **pendingAcceptance**, **accepted**, **rejected** and **failed**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connections` - The list of VPC endpoint service connections.

The `connections` block supports:

* `endpoint_id` - The ID of VPC endpoint.

* `status` - The connection status of the VPC endpoint.

* `domain_id` - The Domain ID.

* `created_at` - The creation time of VPC endpoint.

* `updated_at` - The latest update time of VPC endpoint.

* `description` - The description of the VPC endpoint connection.
