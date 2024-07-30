---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_resource_groups"
description: |-
  Use this data source to get the list of CES resource groups.
---

# huaweicloud_ces_resource_groups

Use this data source to get the list of CES resource groups.

## Example Usage

```hcl
variable "enterprise_project_id " {}
variable "group_name" {}

data "huaweicloud_ces_resource_groups" "test" {
  enterprise_project_id = var.enterprise_project_id
  group_name            = var.group_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the resource group belongs.

* `group_name` - (Optional, String) Specifies the name of a resource group.
  Fuzzy search is supported.

* `group_id` - (Optional, String) Specifies the resource group ID.

* `type` - (Optional, String) Specifies the method of adding resources to a resource group.
  The valid values are as follows:
  + **EPS**: Resources in an enterprise project are added to a resource group.
  + **TAG**: Resources with selected tags are added to a resource group.
  + **Manual**: Resources are added manually to a resource group.
  
  If this parameter is empty, all resource groups are queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resource_groups` - The resource group list.

  The [resource_groups](#resource_groups_struct) structure is documented below.

<a name="resource_groups_struct"></a>
The `resource_groups` block supports:

* `type` - The method of adding resources to a resource group.

* `group_name` - The name of a resource group.

* `group_id` - The resource group ID.

* `created_at` - The time when the resource group was created.

* `enterprise_project_id` - The ID of the enterprise project to which the resource group belongs.
