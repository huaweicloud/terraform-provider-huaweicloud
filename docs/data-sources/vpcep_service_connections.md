---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_service_connections"
description: |-
  Use this data source to get a list of the VPC endpoint service connections.
---

# huaweicloud_vpcep_service_connections

Use this data source to get a list of the VPC endpoint service connections.

## Example Usage

```hcl
variable service_id {}

data "huaweicloud_vpcep_service_connections" "test" {
  service_id  = var.service_id
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to querry the VPC endpoint service connections.
  If omitted, the provider-level region will be used.

* `service_id` - (Required, String) Specifies the ID of VPC endpoint service.

* `endpoint_id` - (Optional, String) Specifies the ID of VPC endpoint which has connected to
  VPC endpoint service.

* `marker_id` - (Optional, String) Specifies the packet ID of the VPC endpoint.

* `status` - (Optional, String) Specifies the connection status of the VPC endpoint.
  The value can be **pendingAcceptance**, **accepted**, **rejected** and **failed**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connections` - The list of VPC endpoint service connections.

The `connections` block supports:

* `endpoint_id` - The ID of VPC endpoint.

* `marker_id` - The packet ID of the VPC endpoint.

* `status` - The connection status of the VPC endpoint.

* `domain_id` - The Domain ID.

* `created_at` - The creation time of VPC endpoint.

* `updated_at` - The latest update time of VPC endpoint.

* `description` - The description of the VPC endpoint connection.
