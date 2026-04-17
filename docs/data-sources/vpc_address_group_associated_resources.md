---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_address_group_associated_resources"
description: |-
  Use this data source to get the list of resources associated with the vpc address group within HuaweiCloud.
---

# huaweicloud_vpc_address_group_associated_resources

Use this data source to get the list of resources associated with the vpc address group within HuaweiCloud.

## Example Usage

```hcl
variable "ip_address_group_id" {}
variable "enterprise_project_id" {}

data "huaweicloud_vpc_address_group_associated_resources" "test" {
  ip_address_group_id   = var.ip_address_group_id
  enterprise_project_id = var.enterprise_project_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `ip_address_group_id` - (Optional, String) Specifies the ID of an IP address group that will be used as a filter.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project
  that an IP address group belongs to that will be used as a filter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `address_groups` - The response body for querying the resources associated with the IP address group.

  The [address_groups](#address_groups_struct) structure is documented below.

<a name="address_groups_struct"></a>
The `address_groups` block supports:

* `id` - The ID of the IP address group.

* `enterprise_project_id` - The ID of the enterprise project that the IP address group belongs to.

* `dependency` - The resources associated with the IP address group.

  The [dependency](#dependency_struct) structure is documented below.

<a name="dependency_struct"></a>
The `dependency` block supports:

* `type` - The type of the resource associated with the IP address group.

* `instance_id` - The ID of the resource associated with the IP address group.

* `instance_name` - The name of the resource associated with the IP address group.
