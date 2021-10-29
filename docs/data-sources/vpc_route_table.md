---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud_vpc_route_table

Provides details about a specific VPC route table.

## Example Usage

```hcl
variable "vpc_id" {}

# get the default route table
data "huaweicloud_vpc_route_table" "default" {
  vpc_id = var.vpc_id
}

# get a custom route table
data "huaweicloud_vpc_route_table" "custom" {
  vpc_id = var.vpc_id
  name   = "demo"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the vpc route table.
  If omitted, the provider-level region will be used.

* `vpc_id` (Required, String) - Specifies the VPC ID where the route table resides.

* `name` (Optional, String) - Specifies the name of the route table.

* `id` (Optional, String) - Specifies the ID of the route table.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `default` (Bool) - Whether the route table is default or not.

* `description` (String) - The supplementary information about the route table.

* `subnets` (List) - An array of one or more subnets associating with the route table.

* `route` (List) - The route object list. The [route object](#route_object) is documented below.

<a name="route_object"></a>
The `route` block supports:

* `type` (String) - The route type.
* `destination` (String) - The destination address in the CIDR notation format
* `nexthop` (String) - The next hop.
* `description` (String) - The description about the route.
