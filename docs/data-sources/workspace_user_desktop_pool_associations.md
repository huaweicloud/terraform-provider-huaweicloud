---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_user_desktop_pool_associations"
description: |-
  Use this data source to query the desktop pools associated with users within HuaweiCloud.
---

# huaweicloud_workspace_user_desktop_pool_associations

Use this data source to query the desktop pools associated with users within HuaweiCloud.

## Example Usage

```hcl
variable "user_ids" {
  type = list(string)
}

data "huaweicloud_workspace_user_desktop_pool_associations" "test" {
  user_ids = var.user_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the user desktop pool associations are located.  
  If omitted, the provider-level region will be used.

* `user_ids` - (Required, List) Specifies the list of user IDs to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `associations` - The list of user associations with desktop pools.  
  The [associations](#workspace_user_desktop_pool_associations) structure is documented below.

<a name="workspace_user_desktop_pool_associations"></a>
The `associations` block supports:

* `user_id` - The ID of the user.

* `desktop_pools` - The list of desktop pools associated with the user.  
  The [desktop_pools](#workspace_user_desktop_pool_associations_desktop_pools) structure is documented below.

<a name="workspace_user_desktop_pool_associations_desktop_pools"></a>
The `desktop_pools` block supports:

* `id` - The ID of the desktop pool.

* `name` - The name of the desktop pool.

* `is_attached` - Whether a desktop is assigned.
