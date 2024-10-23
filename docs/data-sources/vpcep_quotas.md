---
subcategory: "VPC Endpoint (VPCEP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpcep_quotas"
description: |-
  Use this data source to get a list of the VPC endpoint resource quotas.
---

# huaweicloud_vpcep_quotas

Use this data source to get a list of the VPC endpoint resource quotas.

## Example Usage

```hcl
variable type {}

data "huaweicloud_vpcep_quotas" "test" {
  type = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the resource type.
  The value can be **endpoint_service** or **endpoint**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - The list of the VPC endpoint resource quotas.

  The [quotas](#quotas_quotas_struct) structure is documented below.

<a name="quotas_quotas_struct"></a>
The `quotas` block supports:

* `type` - The resource type.

* `used` - The number of used quotas.

* `quota` - The number of available quotas.
