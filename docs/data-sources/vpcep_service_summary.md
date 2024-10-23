---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_service_summary"
description: |-
  Use this data source to get a VPC endpoint service summary information.
---

# huaweicloud_vpcep_service_summary

Use this data source to get a VPC endpoint service summary information.

-> This data source allows current user query the summary information about VPC endpoint service created even if
  it belongs to other users or the system.

## Example Usage

```hcl
variable "endpoint_service_name" {}

data "huaweicloud_vpcep_service_summary" "test" {
  endpoint_service_name = var.endpoint_service_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `endpoint_service_id` - (Optional, String) Specifies the ID of the VPC endpoint service.

* `endpoint_service_name` - (Optional, String) Specifies the name of the VPC endpoint service.

-> Exactly one of `endpoint_service_id` or `endpoint_service_name` must be set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the VPC endpoint service, also is the data source ID.

* `service_name` - The name of the VPC endpoint service.

* `service_type` - The type of the VPC endpoint service.

* `is_charge` - Whether the VPC endpoint connected to the VPC endpoint service is charged.

* `enable_policy` - Whether the VPC endpoint policy can be customized.

* `public_border_group` - The public border group information about the pool corresponding to the VPC endpoint.

* `created_at` - The creation time of the VPC endpoint service, in RFC3339 format.
