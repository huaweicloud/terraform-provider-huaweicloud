---
subcategory: "Resource Access Manager (RAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ram_resource_types"
description: |-
  Use this data source to get the list of resource types in Resource Access Manager.
---

# huaweicloud_ram_resource_types

Use this data source to get the list of resource types in Resource Access Manager.

## Example Usage

```hcl
data "huaweicloud_ram_resource_types" "test" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resource_types` - The list of resource types.

  The [resource_types](#resource_types_struct) structure is documented below.

<a name="resource_types_struct"></a>
The `resource_types` block supports:

* `region_id` - The ID of region.

* `resource_type` - The resource type name.
