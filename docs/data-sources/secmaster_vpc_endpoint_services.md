---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_vpc_endpoint_services"
description: |-
  Use this data source to get the list of SecMaster VPC endpoint services within HuaweiCloud.
---

# huaweicloud_secmaster_vpc_endpoint_services

Use this data source to get the list of SecMaster VPC endpoint services within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_vpc_endpoint_services" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of VPC endpoint services.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `id` - The unique identifier of the service (UUID).

* `name` - The name of the tenant to which the service belongs.

* `type` - The type of VPC service. Valid values are:
  + **MANAGE**: Management channel.
  + **DATA**: Data channel.

* `deprecated` - Whether the service is deprecated.
