---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_private_transit_subnets"
description: |-
  Use this data source to get the list of transit subnets.
---

# huaweicloud_nat_private_transit_subnets

Use this data source to get the list of transit subnets.

## Example Usage

```hcl
data "huaweicloud_nat_private_transit_subnets" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the transit subnets are located.
  If omitted, the provider-level region will be used.

* `ids` - (Optional, List) Specifies the ID of the transit subnet.

* `names` - (Optional, List) Specifies the name of the transit subnet.

* `descriptions` - (Optional, List) Specifies the description of the transit subnet.

* `virsubnet_project_ids` - (Optional, List) Specifies the project ID to which the transit subnet belongs.

* `vpc_ids` - (Optional, List) Specifies the VPC ID to which the transit subnet belongs.

* `virsubnet_ids` - (Optional, List) Specifies the subnet ID to which the transit subnet belongs.

* `status` - (Optional, List) Specifies the status of the transit subnet.
  The value can be **ACTIVE** or **INACTIVE**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `transit_subnets` - The list ot the transit subnets.
  The [transit_subnets](#transit_subnets_struct) structure is documented below.

<a name="transit_subnets_struct"></a>
The `transit_subnets` block supports:

* `id` - The ID of the transit subnet.

* `name` - The name of the transit subnet.

* `description` - The description of the transit subnet.

* `virsubnet_project_id` - The ID of the project to which the transit subnet belongs.

* `project_id` - The project ID.

* `vpc_id` - The ID of the VPC to which the transit subnet belongs.

* `virsubnet_id` - The ID of the subnet to which the transit subnet belongs.

* `cidr` - The CIDR block of the transit subnet.

* `type` - The type of the transit subnet. The value only can be **VPC**.

* `status` - The status of the transit subnet.

* `ip_count` - The number of IP addresses that has been assigned from the transit subnet.

* `created_at` - The creation time of the transit subnet.

* `updated_at` - The latest update time of the transit subnet.

* `tags` - The key/value pairs to associate with the transit subnet.
