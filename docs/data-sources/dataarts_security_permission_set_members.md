---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_permission_set_members"
description: |-
  Use this data source to get the list of associated members under a specified permission set for DataArts Studio
  Security within HuaweiCloud.
---

# huaweicloud_dataarts_security_permission_set_members

Use this data source to get the list of associated members under a specified permission set for DataArts Studio Security
within HuaweiCloud.

## Example Usage

### Query all associated members under a specified permission set

```hcl
variable "workspace_id" {}
variable "permission_set_id" {}

data "huaweicloud_dataarts_security_permission_set_members" "test" {
  workspace_id      = var.workspace_id
  permission_set_id = var.permission_set_id
}
```

### Query the associated members under a specified permission set and using member type to filter

```hcl
variable "workspace_id" {}
variable "permission_set_id" {}

data "huaweicloud_dataarts_security_permission_set_members" "test" {
  workspace_id      = var.workspace_id
  permission_set_id = var.permission_set_id
  member_type       = "USER"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the permission set associated members are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the permission set belongs.

* `permission_set_id` - (Required, String) Specifies the ID of the permission set which the members associated.

* `member_name` - (Optional, String) Specifies the name of the specified permission set member to be queried.

* `member_type` - (Optional, String) Specifies the type of the permission set members to be queried.  
  The valid values are as follows:
  + **USER**
  + **USER_GROUP**
  + **WORKSPACE_ROLE**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `members` - The list of permission set associated members that match the filter parameters.  
  The [members](#dataarts_security_permission_set_associated_members_attr) structure is documented below.

<a name="dataarts_security_permission_set_associated_members_attr"></a>
The `members` block supports:

* `id` - The ID of the permission set member.

* `type` - The type of the permission set member.

* `name` - The name of the permission set member.

* `status` - The status of the permission set member.

* `instance_id` - The instance ID to which the permission set member belongs.

* `workspace` - The workspace ID to which the permission set member belongs.

* `permission_set_id` - The ID of the permission set to which the permission set member belongs.

* `cluster_id` - The cluster ID to which the permission set member belongs.

* `cluster_type` - The cluster type to which the permission set member belongs.

* `cluster_name` - The cluster name to which the permission set member belongs.

* `create_user` - The creator of the permission set member.

* `created_at` - The creation time of the permission set member, in RFC3339 format.

* `deadline` - The deadline time of the permission set member, in RFC3339 format.
