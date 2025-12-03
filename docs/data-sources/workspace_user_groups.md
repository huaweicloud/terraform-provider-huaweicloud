---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_user_groups"
description: |-
  Use this data source to query the Workspace user groups under a specified region within HuaweiCloud.
---

# huaweicloud_workspace_user_groups

Use this data source to query the Workspace user groups under a specified region within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_workspace_user_groups" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the user groups are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - The list of user groups that matched filter parameters.  
  The [groups](#workspace_user_groups) structure is documented below.

<a name="workspace_user_groups"></a>
The `groups` block supports:

* `id` - The ID of the user group.

* `name` - The name of the user group.

* `create_time` - The creation time of the user group, in RFC3339 format.

* `description` - The description of the user group.

* `user_quantity` - The number of users in the user list.

* `parent` - The parent user group of the user group.  
  The [parent](#workspace_user_groups_parent) structure is documented below.

* `realm_id` - The domain ID of the user group.

* `platform_type` - The type of the user group.  
  The valid values are as follows:
  + **AD**: AD domain user group
  + **LOCAL**: Local liteAs user group

* `group_dn` - The distinguished name of the user group.

* `domain` - The domain name of the user group.

* `sid` - The SID of the user group.

* `total_desktops` - The number of users in the user list.

<a name="workspace_user_groups_parent"></a>
The `parent` block supports:

* `id` - The ID of the parent user group.

* `name` - The name of the parent user group.

* `create_time` - The creation time of the parent user group, in RFC3339 format.

* `description` - The description of the parent user group.

* `user_quantity` - The number of users in the parent user group.

* `realm_id` - The domain ID of the parent user group.

* `platform_type` - The type of the parent user group.

* `group_dn` - The distinguished name of the parent user group.

* `domain` - The domain name of the parent user group.

* `sid` - The SID of the parent user group.

* `total_desktops` - The number of users in the parent user group.
