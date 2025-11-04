---
subcategory: "Elastic Cloud Server (ECS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_compute_flavor_capacity"
description: |-
  Use this data source to get the capacity of a flavor.
---

# huaweicloud_compute_flavor_capacity

Use this data source to get the capacity of a flavor.

## Example Usage

```hcl
variable "flavor_id" {}

data "huaweicloud_compute_flavor_capacity" "test" {
  flavor_id = var.flavor_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `flavor_id` - (Required, String) Specifies the flavor ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resources` - Indicates the capacity list of the flavor.

  The [resources](#resources_struct) structure is documented below.

<a name="resources_struct"></a>
The `resources` block supports:

* `region_id` - Indicates the ID of the region that the flavor belongs to.

* `availability_zone` - Indicates the AZ that the flavor belongs to.

* `prefer` - Indicates whether the resources of the flavor in the current AZ are sufficient.
  Value options:
  + **true**: Sufficient
  + **false**: Insufficient
