---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_groups"
description: |-
  Use this data source to get the list of application groups within HuaweiCloud.
---

# huaweicloud_workspace_app_groups

Use this data source to get the list of application groups within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_workspace_app_groups" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `server_group_id` - (Optional, String) Specifies the server group ID associated with the application group.

* `group_id` - (Optional, String) Specifies the ID of the application group.

* `name` - (Optional, String) Specifies the name of the application group.
  Fuzzy search is supported.

* `type` - (Optional, String) Specifies the type of the application group.  
  The valid values are as follows:
  + **COMMON_APP**
  + **SESSION_DESKTOP_APP**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - All application groups that match the filter parameters.

  The [groups](#applications_groups) structure is documented below.

<a name="applications_groups"></a>
The `groups` block supports:

* `id` - The ID of the application group.

* `name` - The name of the application group.

* `server_group_id` - The server group ID associated with the application group.

* `server_group_name` - The server group name associated with the application group.

* `description` - The description of the application group.

* `type` - The type of the application group.

* `app_count` - The number of associated applications.

* `created_at` - The creation time of the application group, in RFC3339 format.
