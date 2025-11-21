---
subcategory: "NAT Gateway (NAT)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_nat_private_transit_subnets"
description: |-
  Use this data source to get the list of transit Subnets.
---

# huaweicloud_nat_private_transit_subnets

Use this data source to get the list of transit subnets.

## Example Usage

```hcl
variable "target_id1" {}
variable "target_id2" {}

data "huaweicloud_nat_private_transit_subnets" "test" {
  ids = [target_id1,target_id2]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the transit subnets are located.
  If omitted, the provider-level region will be used.

* `ids` - (Optional, List) Specifies the resource IDs for querying instances.

* `names` - (Optional, List) Specifies the resource names for querying instances.

* `descriptions` - (Optional, List) Specifies the resource descriptions for querying instances.

* `virsubnet_project_ids` - (Optional, List) Specifies the resource subnet project ids for querying instances.

* `vpc_ids` - (Optional, List) Specifies the resource vpc ids for querying instances.

* `virsubnet_ids` - (Optional, List) Specifies the resource subnet ids for querying instances.

* `status` - (Optional, List) Specifies the resource status for querying instances.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `transit_subnets` - The list ot the transit subnets.
  The [transit_subnets](#private_transitSubnets) structure is documented below.

<a name="private_transitSubnets"></a>
The `transit_subnets` block supports:

* `id` - The ID of the transit subnet.

* `name` - The name of the transit subnet.

* `description` - The description of the transit subnet.

* `virsubnet_project_id` - ID of the project to which the transit subnet belongs.

* `project_id` - The project ID.

* `vpc_id` - ID of the VPC to which the transit subnet belongs.

* `virsubnet_id` - The ID of the subnet to which the transit subnet belongs.

* `cidr` - The CIDR block of the transit subnet.

* `type` - transit subnet type. The value can only be VPC.

* `status` - The status of the transit subnet.

* `ip_count` - The number of IP addresses that has been assigned from the transit subnet.

* `created_at` - The creation time of the transit subnet.

* `updated_at` - The latest update time of the transit subnet.

* `tags` - The tag list.
  The [tags](#nat_resources_tags) structure is documented below.

<a name="nat_resources_tags"></a>
The `tags` block supports:

* `key` - The tag key. Maximum length of 128 Unicode characters. The key cannot be empty.

* `value` - The tag value. Each value has a maximum length of 255 Unicode characters.
