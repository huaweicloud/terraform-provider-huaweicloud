---
subcategory: "Resource Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rc_resource_groups"
description: |-
  Use this data source to get the resource groups within HuaweiCloud Resource Center.
---

# huaweicloud_rc_resource_groups

Use this data source to get the resource groups within HuaweiCloud Resource Center.

## Example Usage

```hcl
data "huaweicloud_rc_resource_groups" "test" {}
```

## Arguments Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `resource_groups` - List of resource groups.
  The [resource_groups](#resource_groups_struct) structure is documented below.

<a name="resource_groups_struct"></a>
The `resource_groups` block supports:

* `group_id` - The resource group ID.

* `domain_id` - The omain ID of the resource group.

* `group_name` - The resource group name.

* `description` - The resource group description.

* `create_time` - The creation time of the resource group.

* `update_time` - The update time of the resource group.
