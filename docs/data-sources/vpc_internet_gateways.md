---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_internet_gateways"
description: ""
---

# huaweicloud_vpc_internet_gateways

Use this data source to get a list of VPC Internet gateways.

## Example Usage

### Example Usage of getting all IGWs

```hcl
data "huaweicloud_vpc_internet_gateways" "all" {}
```

### Example Usage to filter specific IGWs

```hcl
variable "igw_name" {}

data "huaweicloud_vpc_internet_gateways" "filter_by_name" {
  igw_name = var.igw_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `igw_id` - (Optional, String) Specifies the IGW ID.

* `igw_name` - (Optional, String) Specifies the IGW name.

* `vpc_id` - (Optional, String) Specifies the VPC ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `vpc_igws` - The list of VPC IGWs.
  The [vpc_igws](#attrblock--vpc_igws) structure is documented below.

<a name="attrblock--vpc_igws"></a>
The `vpc_igws` block supports:

* `id` - The IGW ID.

* `name` - The IGW name.

* `vpc_id` - The VPC ID to which the IGW associated with.

* `subnet_id` - The subnet ID which the IGW associated with.

* `enable_ipv6` - Indicates the IGW enable ipv6 or not.

* `created_at` - The create time of IGW.

* `updated_at` - The update time of IGW.
